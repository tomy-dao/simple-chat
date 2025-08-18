import { useState, useEffect } from 'react';
import AuthContainer from './components/AuthContainer';
import ChatComponent from './components/Chat';
import auth from './clients/auth';
import { newSocket } from './lib/socket';
import { DefaultEvent } from './lib/socket/socket';


export const socket = newSocket("ws://localhost:8080/chat");

function App() {
  const [user, setUser] = useState(localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')) : null);
  const [isLoading, setIsLoading] = useState(true);
  const [connectId, setConnectId] = useState(null);

  useEffect(() => {
    // Check if user is already logged in
    const token = localStorage.getItem('authToken');
    const savedUser = localStorage.getItem('user');
    
    if (token && savedUser) {
      try {
        setUser(JSON.parse(savedUser));
      } catch (error) {
        console.error('Error parsing saved user:', error);
        localStorage.removeItem('authToken');
        localStorage.removeItem('user');
      }
    }
    
    setIsLoading(false);
  }, []);

  const handleAuthSuccess = (userData) => {
    setUser(userData);
  };

  const handleLogout = async () => {
    await auth.logout();
    setUser(null);
  };

  useEffect(() => {
    if (!user) return;
    socket.connect();
    socket.on(DefaultEvent.Connected, () => {
      socket.emit("authenticate", {
        token: localStorage.getItem("authToken"),
      });
    });
    socket.on("send_connect_id", (connectId) => {
      setConnectId(connectId);
    });
    return () => {
      socket.disconnect();
    };
  }, [user]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="App">
      {user ? (
        <ChatComponent onLogout={handleLogout} user={user} connectId={connectId} />
      ) : (
        <AuthContainer onAuthSuccess={handleAuthSuccess} />
      )}
    </div>
  );
}

export default App;
