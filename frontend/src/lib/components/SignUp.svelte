<script lang="ts">
	import { apiCall } from '$lib/utils/api';
	import { validateField } from '../utils/validate';
	import {
		usernameSchema,
		emailSchema,
		nameSchema,
		passwordWithConfirmationSchema
	} from '../utils/validators';
	import { Eye, EyeOff, ArrowLeft } from 'lucide-svelte';
	import { triggerError, triggerInfo } from '../../stores/modal';
	import { goto } from '$app/navigation';
	let isSubmitting: boolean = false;
	let isChecking: boolean = false;
	let step = 1;
	let showPassword: boolean = false;
	let showConfirmPassword: boolean = false;
	let errorMessage: string = '';
	interface FormData {
		email: string;
		username: string;
		name: string;
		password: string;
		confirmPassword: string;
	}

	let formData: FormData = {
		email: '',
		username: '',
		name: '',
		password: '',
		confirmPassword: ''
	};

	const handleNextButton = async () => {
		errorMessage = '';
		if (step === 1) {
			const emailError = validateField(emailSchema, formData.email);
			if (emailError) {
				errorMessage = emailError;
				return;
			}

			isChecking = true;
			try {
				await apiCall<any>({
					endpoint: `/v1/users/check-email?email=${formData.email}`,
					method: 'GET',
					credentials: 'include',
					onSuccess: () => {
						step++;
					}
				});
			} catch (error) {
				triggerError('An error occurred, please try again later.');
			} finally {
				isChecking = false;
			}
		} else if (step === 2) {
			const usernameError = validateField(usernameSchema, formData.username);
			if (usernameError) {
				errorMessage = usernameError;
				return;
			}

			const nameError = validateField(nameSchema, formData.name);
			if (nameError) {
				errorMessage = nameError;
				return;
			}

			isChecking = true;
			try {
				await apiCall<any>({
					endpoint: `/v1/users/check-username?username=${formData.username}`,
					method: 'GET',
					credentials: 'include',
					onSuccess: () => {
						step++;
					}
				});
			} catch (error) {
				triggerError('An error occurred, please try again later.');
			} finally {
				isChecking = false;
			}
		}
	};
	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && step !== 3) {
			event.preventDefault();
		}
	}

	const handleSubmit = async (): Promise<void> => {
		errorMessage = '';
		const passwordError = validateField(passwordWithConfirmationSchema, {
			password: formData.password,
			confirmPassword: formData.confirmPassword
		});
		if (passwordError) {
			errorMessage = passwordError;
			return;
		}

		isSubmitting = true;
		try {
			await apiCall<any>({
				endpoint: '/v1/auth/register',
				method: 'POST',
				data: formData,
				credentials: 'include',
				onSuccess: () => {
					triggerInfo('Account created successfully!');
					goto('/');
				}
			});
		} catch (error) {
			triggerError('An error occurred, please try again later.');
		} finally {
			isSubmitting = false;
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
			<h1 class="text-4xl font-bold text-indigo-500">BuzzChat</h1>
			<p class="mt-1 text-sm text-gray-400">Create your account</p>
		</div>

		<form on:submit|preventDefault={handleSubmit} on:keydown={handleKeyDown} class="space-y-4">
			{#if step === 1}
				<div>
					<label for="email" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300"
						>Email</label
					>
					<input
						type="text"
						id="email"
						bind:value={formData.email}
						placeholder="Email"
						class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
						required
					/>
				</div>
			{:else if step == 2}
				<div>
					<label for="username" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300"
						>Username</label
					>
					<input
						type="text"
						id="username"
						bind:value={formData.username}
						placeholder="Username"
						class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
						required
					/>
				</div>

				<div>
					<label for="name" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300">Name</label>
					<input
						type="text"
						id="name"
						bind:value={formData.name}
						placeholder="Name"
						class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
					/>
				</div>
			{:else}
				<div>
					<label for="password" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300"
						>Password</label
					>
					<div class="relative">
						<input
							type={showPassword ? 'text' : 'password'}
							id="password"
							bind:value={formData.password}
							placeholder="Password"
							class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 pr-12 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
							required
						/>
						<button
							type="button"
							on:click={() => (showPassword = !showPassword)}
							class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 transition-colors duration-200 hover:cursor-pointer hover:text-indigo-400 focus:outline-none"
							aria-label={showPassword ? 'Hide password' : 'Show password'}
						>
							{#if showPassword}
								<EyeOff class="h-5 w-5" />
							{:else}
								<Eye class="h-5 w-5" />
							{/if}
						</button>
					</div>
				</div>

				<div>
					<label for="confirmPassword" class="block pb-0.5 pl-1 text-sm font-medium text-gray-300"
						>Confirm Password</label
					>
					<div class="relative">
						<input
							type={showConfirmPassword ? 'text' : 'password'}
							id="confirmPassword"
							bind:value={formData.confirmPassword}
							placeholder="Confirm Password"
							class="w-full rounded-xl border border-gray-600 bg-gray-700 px-4 py-2 pr-12 text-white shadow-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-400 focus:outline-none"
							required
						/>
						<button
							type="button"
							on:click={() => (showConfirmPassword = !showConfirmPassword)}
							class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 transition-colors duration-200 hover:cursor-pointer hover:text-indigo-400 focus:outline-none"
							aria-label={showPassword ? 'Hide password' : 'Show password'}
						>
							{#if showConfirmPassword}
								<EyeOff class="h-5 w-5" />
							{:else}
								<Eye class="h-5 w-5" />
							{/if}
						</button>
					</div>
				</div>
			{/if}

			{#if errorMessage !== ''}
				<span class="ml-2 block text-sm text-red-500">{errorMessage}</span>
			{/if}
			{#if step !== 3}
				<button
					type="button"
					on:click={handleNextButton}
					disabled={isChecking}
					class="group relative flex w-full items-center justify-center gap-2 rounded-xl bg-indigo-600 px-4 py-2 font-medium text-white transition hover:cursor-pointer hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 disabled:cursor-default disabled:bg-gray-600"
				>
					{#if isChecking}
						<span
							class="h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
						></span>
						<span></span>
					{:else}
						<span>Next</span>
					{/if}
				</button>
			{:else}
				<button
					type="submit"
					disabled={isSubmitting}
					class="group relative flex w-full items-center justify-center gap-2 rounded-xl bg-indigo-600 px-4 py-2 font-medium text-white transition hover:cursor-pointer hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 disabled:cursor-default disabled:bg-gray-600"
				>
					{#if isSubmitting}
						<span
							class="h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
						></span>
						<span></span>
					{:else}
						<span>Submit</span>
					{/if}
				</button>
			{/if}
		</form>

		<div class="mt-6 text-center text-sm text-gray-400">
			Already have an account?
			<a href="/" class="text-indigo-500 hover:underline">Login</a>
		</div>
	</div>
</main>
