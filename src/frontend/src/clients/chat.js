import { client } from "./initial";

const chat = {
  getConversations: async () => {
    const response = await client.get('/conversations');
    return response.data;
  },
  createConversationByUserId: async (userId) => {
    const response = await client.post(`/conversations`, {
      user_id: userId,
    });
    return response.data;
  },
  getConversationByUserId: async (userId) => {
    const response = await client.get(`/conversations/user/${userId}`);
    return response.data;
  },
  getMessages: async (conversationId) => {
    const response = await client.get(`/conversations/${conversationId}/messages`);
    return response.data;
  },
  sendMessage: async (conversationId, message) => {
    const response = await client.post(`/conversations/${conversationId}/messages`, { ...message });
    return response.data;
  },
};

export default chat;