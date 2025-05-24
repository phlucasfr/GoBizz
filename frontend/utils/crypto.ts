import nacl from 'tweetnacl';
import { apiConfig } from '@/api/config';

export function decrypt(encrypted: any): string {
    let encryptedData: string;

    if (typeof encrypted === 'string') {
        try {
            const parsed = JSON.parse(encrypted);
            encryptedData = parsed.Data;
        } catch {
            encryptedData = encrypted;
        }
    } else if (typeof encrypted === 'object' && encrypted.Data) {
        encryptedData = encrypted.Data;
    } else {
        throw new Error('Invalid cryptography data format');
    }

    const cleanBase64 = encryptedData.replace(/[^A-Za-z0-9+/=]/g, '');
    const padding = (4 - (cleanBase64.length % 4)) % 4;
    const paddedBase64 = cleanBase64 + '='.repeat(padding);

    const rawData = Uint8Array.from(atob(paddedBase64), c => c.charCodeAt(0));

    if (rawData.length < 24) throw new Error('Ciphertext too short');

    const nonce = rawData.slice(0, 24);
    const ciphertext = rawData.slice(24);

    const keyBytes = new TextEncoder().encode(apiConfig.masterKey);
    const decrypted = nacl.secretbox.open(ciphertext, nonce, keyBytes);

    if (!decrypted) {
        throw new Error('Decrypt error');
    }

    return new TextDecoder().decode(decrypted as Uint8Array);
}


function uint8ArrayToBase64(bytes: Uint8Array): string {
    let binary = '';
    const len = bytes.byteLength;
    for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
}

export function encrypt(data: string): string {
    const keyBytes = new TextEncoder().encode(apiConfig.masterKey);
    const nonce = nacl.randomBytes(24);
    const ciphertext = nacl.secretbox(new TextEncoder().encode(data), nonce, keyBytes);

    const combined = new Uint8Array(nonce.length + ciphertext.length);
    combined.set(nonce);
    combined.set(ciphertext, nonce.length);

    return uint8ArrayToBase64(combined);
}
