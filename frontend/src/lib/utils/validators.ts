import { z } from "zod";

export const usernameSchema = z
    .string()
    .max(30, "Username must be less than 30 characters.")
    .regex(/^[a-zA-Z0-9_]+$/, "Username should only contain alphanumeric characters and underscores")
    .nonempty("Username is required!");

export const passwordSchema = z
    .string()
    .min(8, "Password must be at least 8 characters.")
    .nonempty("Password is required!");

export const emailSchema = z
    .string()
    .regex(
        /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
        "Enter a valid email!"
    )
    .nonempty("Email is required!");

export const nameSchema = z
    .string()
    .max(50, "Name must be less than 50 characters.")

export const passwordWithConfirmationSchema = z
    .object({
        password: z
            .string()
            .min(8, "Password must be at least 8 characters.")
            .nonempty("Password is required!"),
        confirmPassword: z
            .string()
            .nonempty("Confirm Password is required!")
    })
    .refine((data) => data.password === data.confirmPassword, {
        message: "Passwords do not match!",
        path: ["confirmPassword"],
    });
