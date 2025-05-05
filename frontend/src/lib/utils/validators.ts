import { z } from "zod";

export const usernameSchema = z
    .string()
    .min(1, "Username is required!")
    .max(30, "Useranme must be less than 30 characters");

export const passwordSchema = z
    .string()
    .min(8, "Password must be at least 8 characters");
