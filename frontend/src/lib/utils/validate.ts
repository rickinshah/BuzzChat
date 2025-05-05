import { ZodSchema } from 'zod';

export function validateField<T>(schema: ZodSchema<T>, data: unknown): string | null {
    const result = schema.safeParse(data);

    if (!result.success) {
        return result.error.issues[0].message;
    }

    return null;
}
