import { writable } from "svelte/store";

export const errorMessage = writable();
export const infoMessage = writable();
export const showErrorModal = writable(false);
export const showInfoModal = writable(false);

export function triggerError(message: string | Object) {
	errorMessage.set(message);
	showErrorModal.set(true);
}

export function triggerInfo(message: string | Object) {
	infoMessage.set(message)
	showInfoModal.set(true);
}
