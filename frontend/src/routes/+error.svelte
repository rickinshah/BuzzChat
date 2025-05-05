<script lang="ts">
	import { goto } from '$app/navigation';
	const goHome = () => goto('/');
</script>

<main
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-gray-800 via-gray-900 to-black px-4"
>
	<div
		class="w-full max-w-md rounded-3xl bg-gray-800/90 p-10 text-center shadow-xl backdrop-blur-md transition-all duration-300 hover:shadow-2xl"
	>
		<h1 class="text-7xl font-extrabold text-indigo-500">404</h1>
		<p class="mt-4 text-2xl font-semibold text-white">Page Not Found</p>
		<p class="mt-2 text-gray-400">The page you're looking for doesn't exist or was moved.</p>
		<button
			on:click={goHome}
			class="mt-6 rounded-xl bg-indigo-600 px-5 py-2 text-sm font-medium text-white transition hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
		>
			← Go back home
		</button>
	</div>
</main>
