import { decrypt, encrypt } from '@/utils/crypto';

export class EncryptionService {
    static encryptData(data: unknown): { data: string } {
        const encryptedData = encrypt(JSON.stringify(data));
        return { data: encryptedData };
    }

    static decryptAndParse<T>(data: unknown): T {
        try {
            // If data is already a string, use it directly
            if (typeof data === 'string') {
                const decryptedData = decrypt({ Data: data });
                return JSON.parse(decryptedData);
            }

            // If data is an object with a Data property
            if (typeof data === 'object' && data !== null && 'Data' in data) {
                const decryptedData = decrypt({ Data: (data as { Data: string }).Data });
                return JSON.parse(decryptedData);
            }

            throw new Error('Invalid data format for decryption');
        } catch (error) {
            console.error('Decryption error:', error);
            throw new Error('Failed to process response data');
        }
    }
} 