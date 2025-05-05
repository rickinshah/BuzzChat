<script lang="ts">
	import '../app.css';
	import ErrorModal from '$lib/components/ErrorModal.svelte';
	import { errorMessage, showErrorModal, showInfoModal, infoMessage } from '../stores/modal';
	import InfoModal from '$lib/components/InfoModal.svelte';

	const closeErrorModal = () => {
		showErrorModal.set(false);
		errorMessage.set('');
	};

	const closeInfoModal = () => {
		showInfoModal.set(false);
		infoMessage.set('');
	};

	let { children } = $props();
</script>

{#if $showErrorModal}
	<ErrorModal errorMessage={$errorMessage} onClose={closeErrorModal} />
{/if}

{#if $showInfoModal}
	<InfoModal infoMessage={$infoMessage} onClose={closeInfoModal} />
{/if}

{@render children()}
