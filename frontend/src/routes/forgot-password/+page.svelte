<script lang="ts">
	import { apiCall } from '$lib/utils/api';
	import { validateField } from '$lib/utils/validate';
	import { usernameSchema } from '$lib/utils/validators';
	import { ArrowLeft } from 'lucide-svelte';
	import { triggerError, triggerInfo } from '../../stores/modal';

	let formData = { username: '', otp: '' };
	let errors: Record<string, string> = {};
	let isSubmitting: boolean = false;
	let step: number = 1;

	const handleSubmit = async () => {
		errors = {};

		if (step === 1) {
			const usernameError = validateField(usernameSchema, formData.username);
			if (usernameError) {
				errors.username = usernameError;
				triggerInfo(Object(errors));
				return;
			}

			isSubmitting = true;
			try {
				await apiCall<any>({
					endpoint: '/v1/auth/otp',
					method: 'POST',
					data: { username: formData.username },
					onSuccess: (response) => {
						triggerInfo('OTP sent! Check your registered email');
						step++;
					}
				});
			} catch (error) {
				triggerError('Failed to send OTP. Please try again later.');
			} finally {
				isSubmitting = false;
			}
		} else if (step === 2) {
			isSubmitting = true;
			try {
				await apiCall<any>({
					endpoint: '/v1/auth/otp/validate',
					method: 'POST',
					data: { username: formData.username, otp: formData.otp },
					onSuccess: (response) => {
						triggerInfo('Correct OTP!');
					}
				});
			} catch (error) {
				triggerError('Failed to send OTP. Please try again later.');
			} finally {
				isSubmitting = false;
			}
		}
	};
</script>

<main
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-gray-800 via-gray-900 to-black px-4"
>
	<div
		class="w-full max-w-md rounded-3xl bg-gray-800/90 p-10 shadow-xl backdrop-blur-md transition-all duration-300 hover:shadow-2xl"
	>
		{#if step !== 1}
			<button
				class="fixed flex h-5 w-5 items-center justify-center rounded-full"
				aria-label="Back"
				on:click={() => step--}
			>
				<ArrowLeft
					class="h-5 w-5 text-white transition-all duration-200 ease-in-out hover:cursor-pointer hover:text-indigo-500"
				/>
			</button>
		{/if}
		<div class="mb-6 text-center">
			<h1 class="text-3xl font-bold text-indigo-500">Forgot Password</h1>
			{#if step === 1}
				<p class="mt-1 text-sm text-gray-400">Enter your username to receive a reset code.</p>
			{:else if step === 2}
				<p class="mt-1 text-sm text-gray-400">Enter your OTP to reset password</p>
			{/if}
		</div>

		<form on:submit|preventDefault={handleSubmit} class="space-y-4">
			{#if step === 1}
				<div>
					<label for="username" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300">
						Username/Email
					</label>
					<input
						type="text"
						id="username"
						bind:value={formData.username}
						placeholder="Username/Email"
						class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
						required
					/>
				</div>
			{:else if step === 2}
				<div>
					<label for="username" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300">
						OTP
					</label>
					<input
						type="text"
						id="otp"
						bind:value={formData.otp}
						placeholder="OTP"
						class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
						required
					/>
				</div>
			{/if}

			<button
				type="submit"
				disabled={isSubmitting}
				class="group relative flex w-full items-center justify-center gap-2 rounded-xl bg-indigo-600 px-4 py-2 font-medium text-white transition hover:cursor-pointer hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 disabled:bg-gray-600"
			>
				{#if isSubmitting}
					<span class="h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
					></span>
					<span></span>
				{:else}
					<span>Send OTP</span>
				{/if}
			</button>
		</form>

		<div class="mt-6 text-center text-sm text-gray-400">
			Remembered your password?
			<a href="/" class="text-indigo-500 hover:underline">Back to Login</a>
		</div>
	</div>
</main>
