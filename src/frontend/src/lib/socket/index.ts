import { SocketListener } from './listener';
import { Socket } from './socket';

export function newSocket(url: string) {
    const socket = new Socket(url);
    return socket;
}

export function newSocketListener() {
    return new SocketListener();
}
