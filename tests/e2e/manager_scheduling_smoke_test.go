//go:build e2e

package e2e_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	clientchat "github.com/Pickausernaame/chat-service/tests/e2e/client-chat"
	managerworkspace "github.com/Pickausernaame/chat-service/tests/e2e/manager-workspace"
	wsstream "github.com/Pickausernaame/chat-service/tests/e2e/ws-stream"
)

var _ = Describe("Manager Scheduling Smoke", Ordered, func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc

		clientChat        *clientchat.Chat
		clientStream      *wsstream.Stream
		clientStreamErrCh = make(chan error, 1)

		managerWs          *managerworkspace.Workspace
		managerStream      *wsstream.Stream
		managerStreamErrCh = make(chan error, 1)
	)

	BeforeAll(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Setup client.
		clientChat = newClientChat(ctx, clientsPool.Get())

		var err error
		clientStream, err = wsstream.New(wsstream.NewOptions(
			wsClientEndpoint,
			wsClientOrigin,
			wsClientSecProtocol,
			clientChat.AccessToken(),
			clientChat.HandleEvent,
		))
		Expect(err).ShouldNot(HaveOccurred())
		go func() { clientStreamErrCh <- clientStream.Run(ctx) }()

		// Setup manager.
		managerWs = newManagerWs(ctx, managersPool.Get())

		managerStream, err = wsstream.New(wsstream.NewOptions(
			wsManagerEndpoint,
			wsManagerOrigin,
			wsManagerSecProtocol,
			managerWs.AccessToken(),
			managerWs.HandleEvent,
		))
		Expect(err).ShouldNot(HaveOccurred())
		go func() { managerStreamErrCh <- managerStream.Run(ctx) }()
	})

	AfterAll(func() {
		cancel()
		Expect(<-clientStreamErrCh).ShouldNot(HaveOccurred())
		Expect(<-managerStreamErrCh).ShouldNot(HaveOccurred())
	})

	It("no chats at the start of working day", func() {
		err := managerWs.Refresh(ctx)
		Expect(err).ShouldNot(HaveOccurred())

		n := managerWs.ChatsCount()
		Expect(n).Should(Equal(0))
	})

	It("manager assigned to new problem", func() {
		err := managerWs.ReadyToNewProblems(ctx)
		Expect(err).ShouldNot(HaveOccurred())

		err = clientChat.SendMessage(ctx, "Hello, sir!")
		Expect(err).ShouldNot(HaveOccurred())

		// Client side.

		waitForEvent(clientStream) // NewMessageEvent.
		waitForEvent(clientStream) // MessageSentEvent.
		waitForEvent(clientStream) // NewMessageEvent (service).

		msg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(msg.Body).Should(Equal(fmt.Sprintf("Manager %s will answer you", managerWs.ManagerID())))

		// Manager side.

		waitForEvent(managerStream)         // NewChatEvent.
		waitForOptionalEvent(managerStream) // NewMessageEvent.

		n := managerWs.ChatsCount()
		Expect(n).Should(Equal(1))

		newChat, ok := managerWs.LastChat()
		Expect(ok).Should(BeTrue())
		Expect(newChat.ClientID.String()).Should(Equal(clientChat.ClientID().String()))
		Expect(newChat.ID).ShouldNot(BeEmpty())
	})

	It("assigned problem does not disappear", func() {
		err := managerWs.Refresh(ctx)
		Expect(err).ShouldNot(HaveOccurred())

		n := managerWs.ChatsCount()
		Expect(n).Should(Equal(1))
	})

	It("manager see chat history", func() {
		lastChat, ok := managerWs.LastChat()
		Expect(ok).Should(BeTrue())

		n := lastChat.MessagesCount()
		Expect(n).Should(Equal(1))

		lastMsg, ok := lastChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(lastMsg.ID).ShouldNot(BeEmpty())
		Expect(lastMsg.ChatID).Should(Equal(lastChat.ID))
		Expect(lastMsg.AuthorID.String()).Should(Equal(clientChat.ClientID().String()))
		Expect(lastMsg.CreatedAt.IsZero()).Should(BeFalse())
	})

	It("manager answers back", func() {
		lastChat, ok := managerWs.LastChat()
		Expect(ok).Should(BeTrue())
		err := managerWs.SendMessage(ctx, lastChat.ID)
		Expect(err).ShouldNot(HaveOccurred())

		waitForEvent(clientStream) // NewMessageEvent.

		clientMsg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(clientMsg.AuthorID.String()).Should(Equal(managerWs.ManagerID().String()))
		Expect(clientMsg.Body).Should(Equal("hello"))
		Expect(clientMsg.CreatedAt.IsZero()).Should(BeFalse())

		waitForEvent(managerStream) // NewMessageEvent.
		managerWs.Refresh(ctx)
		lastChat, _ = managerWs.LastChat()

		managerMsg, ok := lastChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(managerMsg.AuthorID.String()).Should(Equal(managerWs.ManagerID().String()))
		Expect(managerMsg.Body).Should(Equal("hello"))
		Expect(managerMsg.CreatedAt.IsZero()).Should(BeFalse())

		GinkgoWriter.Println(lastChat.Messages())
		n := lastChat.MessagesCount()
		Expect(2).Should(Equal(n))
	})

	It("manager closes chat", func() {
		count := managerWs.ChatsCount()
		lastChat, ok := managerWs.LastChat()
		Expect(ok).Should(BeTrue())
		err := managerWs.CloseChat(ctx, lastChat.ID)
		Expect(err).ShouldNot(HaveOccurred())

		waitForEvent(managerStream) // ChatClosedEvent.
		managerWs.Refresh(ctx)
		Expect(managerWs.ChatsCount()).Should(Equal(count - 1))

		Expect(managerWs.CanTakeMoreProblems()).Should(BeTrue())

		waitForEvent(clientStream) // NewMessageEvent
		clientMsg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(clientMsg.Body).Should(Equal("Your question has been marked as resolved.\nThank you for being with us!"))
		Expect(clientMsg.CreatedAt.IsZero()).Should(BeFalse())
		Expect(clientMsg.IsService).Should(BeTrue())

		Expect(4).Should(Equal(clientChat.MessagesCount()))
	})
})
