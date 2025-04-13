'use client';

import { Suspense } from 'react';
import EmailVerificationContent from './EmailVerificationContent';

export default function EmailVerificationPage() {
    return (
        <Suspense fallback={<div>Carregando...</div>}>
            <EmailVerificationContent />
        </Suspense>
    );
}
