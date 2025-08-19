/* eslint-disable @typescript-eslint/no-explicit-any */

export type ListenerCallback = (_message: any, _event: string) => void
type Listener = {[event: string]: ListenerCallback[]}

export class SocketListener {
    listener: Listener = {};

    add(event: string, fn: ListenerCallback) {
        // if (!this) return;
        if (Object.prototype.hasOwnProperty.call(this.listener, event)) {
            this.listener[event].push(fn);
        } else {
            this.listener[event] = [fn];
        }

        return {
            remove: () => {
                this.remove(event, fn);
            },
        }
    }

    get(event: string) {
        if (!Object.prototype.hasOwnProperty.call(this.listener, event)) return [];
        return this.listener[event];
    }

    remove(event: string, fn: ListenerCallback) {
        if (this.listener[event]) {
            const index = this.listener[event].findIndex((cb) => cb === fn);
            if (index !== -1) {
                this.listener[event].splice(index, 1);
            }
        }
    }

    trigger(event: string, payload: any) {
        const listeners = this.get(event);
        listeners.forEach((fn) => {
            fn(payload, event);
        });
    }

    on(event: string, fn: ListenerCallback) {
        return this.add(event, fn);
    }

    emit(event: string, payload: any) {
        this.trigger(event, payload);
    }
}
