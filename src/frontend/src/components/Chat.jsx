import { useState, useEffect } from "react";
import { Chat } from "@/components/ui/chat";
import ChatList from "@/components/ChatList";
import chat from "@/clients/chat";
import { socket, eventListener } from "@/App";

let count = 2;

export function ChatComponent({ onLogout, user, connectId }) {
  const [conversationId, setConversationId] = useState(null);
  const [userId, setUserId] = useState(null);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");

  const handleInputChange = (e) => {
    setInput(e.target.value);
  };

  const sendMessageWhenConversationSelected = async (
    conversationId,
    messageTerminated
  ) => {
    setMessages([...messages, messageTerminated]);
    setInput("");
    const { data: message } = await chat.sendMessage(conversationId, {
      content: messageTerminated.content,
      session_id: connectId,
    });
    eventListener.emit("message", { message });
    setMessages((prevMessages) => {
      // replace the messageTerminated with the new message
      const index = prevMessages.findIndex(
        (message) => message.id === messageTerminated.id
      );
      if (index !== -1) {
        prevMessages[index] = {
          id: message.id,
          role: user.id === message.sender_id ? "user" : "assistant",
          content: message.content,
        };
      }
      return [...prevMessages];
    });
  };

  const sendMessageWhenUserSelected = async (message) => {
    const { data: cvs } = await chat.createConversationByUserId(userId);
    setConversationId(cvs.id);
    sendMessageWhenConversationSelected(cvs.id, message);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!input.trim()) return;

    const messageTerminated = {
      id: `message-terminated-${count++}`,
      role: "user",
      content: input.trim(),
    };
    if (conversationId) {
      sendMessageWhenConversationSelected(conversationId, messageTerminated);
    } else {
      sendMessageWhenUserSelected(messageTerminated);
    }
  };

  const handleConversationSelect = async (conversationId) => {
    setConversationId(conversationId);
    setMessages([]);
    const { data: messages } = await chat.getMessages(conversationId);
    setMessages(
      messages
        .map((message) => ({
          id: message.id,
          role: user.id === message.sender_id ? "user" : "assistant",
          content: message.content,
          // createdAt: new Date(message.created_at),
        }))
        .reverse()
    );
  };

  const handleUserSelect = (userId) => {
    setUserId(userId);
    setMessages([]);
  };

  useEffect(() => {
    if (!conversationId) return;
    const msgSignal = socket.on("message", ({ message }) => {
      if (message.conversation_id == conversationId) {
        setMessages((prevMessages) => [
          ...prevMessages,
          {
            id: message.id,
            role: user.id === message.sender_id ? "user" : "assistant",
            content: message.content,
          },
        ]);
      }
    });

    return () => {
      msgSignal.remove();
    };
  }, [conversationId]);

  // const isLoading = status === "submitted" || status === "streaming"

  return (
    <div
      style={{
        height: "100vh",
        background: "linear-gradient(135deg, #eff6ff 0%, #f8fafc 100%)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: "1rem",
      }}
    >
      <ChatList
        onLogout={onLogout}
        conversationId={conversationId}
        userId={userId}
        user={user}
        onConversationSelect={handleConversationSelect}
        onUserSelect={handleUserSelect}
      />
      {!conversationId?.toString() && !userId?.toString() ? (
        <div className="flex-1 h-full border rounded-lg pt-4 bg-white flex items-center justify-center">
          <p className="text-center text-gray-500 text-2xl">
            Select a conversation to start chatting <br />
            or choose a user to start a new conversation
          </p>
        </div>
      ) : (
        <Chat
          className="flex-1 h-full border rounded-lg pt-4 bg-white"
          messages={messages}
          input={input}
          handleInputChange={handleInputChange}
          handleSubmit={handleSubmit}
          setMessages={setMessages}
        />
      )}
    </div>
  );
}

export default ChatComponent;
