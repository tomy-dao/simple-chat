/* eslint-disable no-console */
import { ListenerCallback, SocketListener } from './listener';
import { newSocketListener } from '.';

/* eslint-disable @typescript-eslint/no-explicit-any */
type SocketMessage = {
    event: string
    payload: any
}

const configKeepConnectInterval = 10000;

export const DefaultEvent = {
    Error: 'error',
    Disconnect: 'disconnect',
    Connected: 'connected',
    Message: 'new_message',
};

export class Socket {
    url = '';

    ws!: WebSocket;

    listener: SocketListener = newSocketListener();

    connecting = false;

    constructor(url = '') {
        this.url = url;
    }

    connect() {
        if (!this.url && this.connecting) return;

        let keepConnectInterval: NodeJS.Timeout;
        const emitKeepConnect = () => this.emit('keep_connect', null);

        this.ws = new WebSocket(this.url);

        this.ws.onmessage = (event) => {
            if (keepConnectInterval) {
                clearInterval(keepConnectInterval);
                keepConnectInterval = setInterval(emitKeepConnect, configKeepConnectInterval);
            }
            const data: SocketMessage = JSON.parse(event.data);

            this.trigger(DefaultEvent.Message, data);

            if (data?.event) {
                this.trigger(data.event, data.payload);
            } else {
                console.log(event);
            }
        };

        this.ws.onclose = (event) => {
            if (keepConnectInterval) clearInterval(keepConnectInterval);

            this.trigger(DefaultEvent.Disconnect, event);
            this.connecting = false;
        };

        this.ws.onopen = (event) => {
            this.connecting = true;
            this.trigger(DefaultEvent.Connected, event);
            keepConnectInterval = setInterval(emitKeepConnect, configKeepConnectInterval);
        };

        this.ws.onerror = (event) => {
            this.trigger(DefaultEvent.Error, event);
            this.connecting = false;
        };
    }

    emit(event: string, payload: any = null) {
        if (!this.connecting) {
            console.error('socket is not connected');
            return;
        }

        this.ws.send(JSON.stringify({
            event,
            payload,
        }));
    }

    onConnected(fn: ListenerCallback) {
        this.on(DefaultEvent.Connected, fn);
    }

    onDisconnected(fn: ListenerCallback) {
        this.on(DefaultEvent.Disconnect, fn);
    }

    onError(fn: ListenerCallback) {
        this.on(DefaultEvent.Error, fn);
    }

    on(event: string, fn: ListenerCallback) {
        this.listener.add(event, fn);
        return {
            remove: () => {
                this.listener.remove(event, fn);
            },
        }
    }

    removeListener(event: string, fn: ListenerCallback) {
        this.listener.remove(event, fn);
    }

    disconnect() {
        if (this.ws) this.ws.close();
    }

    trigger(event: SocketMessage['event'], payload: SocketMessage['payload'] = null) {
        const listeners = this.listener.get(event);
        listeners.forEach((fn) => {
            try {
                fn(payload, event);
            } catch (e) {
                console.log(e);
            }
        });
    }
}
