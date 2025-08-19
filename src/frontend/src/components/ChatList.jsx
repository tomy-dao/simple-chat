import { useState, useEffect } from "react";
import auth from "@/clients/auth";
import chat from "@/clients/chat";
import { socket, eventListener } from "@/App";
import { toast } from "sonner";

const ChatList = ({ onLogout, userId, onUserSelect, onConversationSelect, user, conversationId }) => {
  const [users, setUsers] = useState([]);
  const [conversations, setConversations] = useState([]);
  const [activeTab, setActiveTab] = useState("conversations");

  const handleTabClick = (tab) => {
    setActiveTab(tab);
  };

  const fetchConversations = async () => {
    const { data: conversations } = await chat.getConversations();
    setConversations(conversations);
  };

  const fetchUsers = async () => {
    const { data: users } = await auth.getUsers();
    setUsers(users);
  };

  const handleConversationSelect = async (conversation) => {
    onConversationSelect(conversation.id);
    const userSelected = conversation.user_ids.find(
      (uId) => uId !== user.id
    );
    onUserSelect(userSelected || user.id);
  };

  const handleUserSelect = async (userId) => {
    try {
      const { data: conversation } = await chat.getConversationByUserId(userId);
      onConversationSelect(conversation.id);
      const userIdSelected = conversation.user_ids.find(
        (uId) => uId == userId
      );
      onUserSelect(userIdSelected);
    } catch {
      onUserSelect(userId);
      onConversationSelect(null);
    }
  };

  useEffect(() => {
    if (activeTab === "conversations") {
      fetchConversations();
    } else {
      fetchUsers();
    }
  }, [activeTab]);

  useEffect(() => {
    const msgSignal = socket.on("message", ({message}) => {
      console.log(message);
      if (message.conversation_id !== conversationId && message.sender_id !== user.id) {
        toast.success("New message received", {
          position: "top-right",
          style: {
            background: "green",
            color: "white",
          },
        });
      }
    });
    return () => {
      msgSignal.remove();
    }
  }, [conversationId]);

  useEffect(() => {
    function updateConversation({message}) {
      setConversations((prevConversations) => {
        console.log(message.conversation_id);
        const existingIndex = prevConversations.findIndex(conv => conv.id === message.conversation_id);
        if (existingIndex !== -1) {
          return [message.conversation, ...prevConversations.filter((_, index) => index !== existingIndex)];
        }
        return [message.conversation, ...prevConversations];
      });
    }
    const msgSignal = socket.on("message", updateConversation);
    const msgListener = eventListener.on("message", updateConversation);

    return () => {
      msgSignal.remove();
      msgListener.remove();
    }
  }, []);

  return (
    <div className="w-70 bg-white border rounded-lg mr-4 h-full overflow-y-auto flex flex-col">
      <div className="p-4 border-b">
        <div className="flex border-b">
          <button
            className={`flex-1 px-4 py-2 text-sm font-medium text-gray-600 border-b-2 ${
              activeTab === "conversations"
                ? "border-blue-500 text-blue-600"
                : "hover:text-gray-800"
            }`}
            onClick={() => handleTabClick("conversations")}
          >
            Conversations
          </button>
          <button
            className={`flex-1 px-4 py-2 text-sm font-medium text-gray-600 border-b-2 ${
              activeTab === "people"
                ? "border-blue-500 text-blue-600"
                : "hover:text-gray-800"
            }`}
            onClick={() => handleTabClick("people")}
          >
            People
          </button>
        </div>
      </div>
      <div className="flex flex-col h-full">
        <div className="p-2">
          <div className="space-y-2">
            {activeTab === "people" ? (
              users.map((user) => (
                <div
                  key={user.id}
                  className={`p-3 hover:bg-gray-50 rounded-lg cursor-pointer border-l-4 bg-blue-50 ${user.id === userId ? "border-blue-500" : ""}`}
                  onClick={() => handleUserSelect(user.id)}
                >
                  <div className="font-medium text-sm text-gray-800">
                    {user.username}
                  </div>
                </div>
              ))
            ) : <></>}
            {console.log(conversations)}
            {activeTab === "conversations" ? (
              conversations.map((conversation) => (
                <div
                  key={conversation.id}
                  className={`p-3 hover:bg-gray-50 rounded-lg cursor-pointer border-l-4 bg-blue-50 ${conversation?.id === conversationId ? "border-blue-500" : ""}`}
                  onClick={() => handleConversationSelect(conversation)}
                >
                  <div className="font-medium text-sm text-gray-800">
                    {conversation.participants.find(
                      (participant) => participant.user_id !== user.id
                    )?.user?.username ?? conversation.participants[0].user.username}
                  </div>
                </div>
              ))
            ) : <></>}
          </div>
        </div>
      </div>
      <div className="p-4 border-t mt-auto">
        <button
          className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm"
          onClick={onLogout}
        >
          Logout
        </button>
      </div>
    </div>
  );
};

export default ChatList;
