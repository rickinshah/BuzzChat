import { writable } from "svelte/store";

export const goto = writable<string | null>(null);
