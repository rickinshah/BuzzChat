<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { scale } from 'svelte/transition';
	export let errorMessage: string | Object;
	export let onClose: () => void;

	let previousActive: HTMLElement | null = null;

	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') {
			onClose();
		}
	};

	onMount(() => {
		window.addEventListener('keydown', handleKeydown);
		previousActive = document.activeElement as HTMLElement;
		previousActive?.blur();
	});

	onDestroy(() => {
		window.removeEventListener('keydown', handleKeydown);
		previousActive?.focus();
	});
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
	<div
		transition:scale={{ duration: 200 }}
		class="relative w-full max-w-md rounded-2xl border border-white/10 bg-[#1c1c1e]/80 p-6 text-white shadow-2xl backdrop-blur-md"
	>
		<!-- Icon -->
		<div class="mb-2 flex items-center justify-center">
			<div class="flex h-12 w-12 items-center justify-center rounded-full bg-red-500/10">
				<svg
					class="h-6 w-6 text-red-500"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</div>
		</div>

		<!-- Message -->
		<h2 class="mb-1 text-center text-xl font-semibold text-red-400">Error</h2>
		{#if typeof errorMessage === 'string'}
			<p class="mb-4 text-center text-sm break-words text-gray-300">{errorMessage}</p>
		{:else}
			<ul class="mb-4 text-center text-sm break-words text-gray-300">
				{#each Object.entries(errorMessage) as [field, message]}
					<li><span class="font-medium capitalize">{field}:</span> {message}</li>
				{/each}
			</ul>
		{/if}

		<!-- Action -->
		<div class="text-center">
			<button
				on:click={onClose}
				class="inline-flex cursor-pointer items-center justify-center rounded-xl bg-indigo-600 px-5 py-2 text-sm font-medium text-white transition hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
			>
				Close
			</button>
		</div>

		<!-- Close button (top-right) -->
		<button
			on:click={onClose}
			class="absolute top-3 right-3 cursor-pointer text-xl text-gray-400 hover:text-white"
			aria-label="Close"
		>
			&times;
		</button>
	</div>
</div>
