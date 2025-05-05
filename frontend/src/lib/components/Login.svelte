<script lang="ts">
	import { apiCall } from '$lib/utils/api';
	import ErrorModal from './ErrorModal.svelte';
	import { validateField } from '../utils/validate';
	import { passwordSchema, usernameSchema } from '../utils/validators';
	import { triggerError, triggerInfo } from '../../stores/modal';
	let showErrorModal: boolean = false;
	let isSubmitting: boolean = false;
	let errors: Record<string, string> = {};
	interface FormData {
		username: string;
		password: string;
	}

	let formData: FormData = {
		username: '',
		password: ''
	};

	const handleSubmit = async (): Promise<void> => {
		errors = {};

		const usernameError = validateField(usernameSchema, formData.username);
		if (usernameError) {
			errors.username = usernameError;
		}
		const passwordError = validateField(passwordSchema, formData.password);
		if (passwordError) {
			errors.password = passwordError;
		}

		if (Object.keys(errors).length > 0) {
			triggerError(Object(errors));
			return;
		}

		isSubmitting = true;
		try {
			await apiCall<any>({
				endpoint: '/v1/auth/login',
				method: 'POST',
				data: formData,
				onSuccess: (response) => {
					localStorage.setItem('user', JSON.stringify(response.user));
					triggerInfo('Login Successful');
				}
			});
		} catch (error) {
			triggerError('An error occurred, please try again later.');
		} finally {
			isSubmitting = false;
		}
	};
	const closeModal = () => {
		showErrorModal = false;
	};
</script>

<main
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-gray-800 via-gray-900 to-black px-4"
>
	<div
		class="w-full max-w-md rounded-3xl bg-gray-800/90 p-10 shadow-xl backdrop-blur-md transition-all duration-300 hover:shadow-2xl"
	>
		<div class="mb-6 text-center">
			<h1 class="text-4xl font-bold text-indigo-500">BuzzChat</h1>
			<p class="mt-1 text-sm text-gray-400">Welcome back! Log in to start messaging.</p>
		</div>

		<form on:submit|preventDefault={handleSubmit} class="space-y-4">
			<div>
				<label for="username" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300"
					>Username/Email</label
				>
				<input
					type="text"
					id="username"
					bind:value={formData.username}
					placeholder="Username/Email"
					class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
					required
				/>
			</div>

			<div>
				<label for="password" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300"
					>Password</label
				>
				<input
					type="password"
					id="password"
					bind:value={formData.password}
					placeholder="Password"
					class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
					required
				/>
			</div>

			<!-- Forgot Password Link Above Login Button -->
			<div class="text-right text-sm">
				<a href="/forgot-password" class="text-indigo-400 hover:underline">Forgot your password?</a>
			</div>

			<button
				type="submit"
				disabled={isSubmitting}
				class="group relative flex w-full items-center justify-center gap-2 rounded-xl bg-indigo-600 px-4 py-2 font-medium text-white transition hover:cursor-pointer hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 disabled:bg-gray-600"
			>
				{#if isSubmitting}
					<span class="h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
					></span>
					<span>Logging in...</span>
				{:else}
					<span>Login</span>
				{/if}
			</button>
		</form>

		<div class="mt-6 text-center text-sm text-gray-400">
			Don't have an account?
			<a href="/signup" class="text-indigo-500 hover:underline">Sign up</a>
		</div>
	</div>
</main>
