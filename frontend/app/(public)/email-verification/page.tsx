'use client';

import { Button } from '@/components/ui/button';
import { Loader2 } from 'lucide-react';
import { apiConfig } from '@/api/config';
import { apiRequest } from '@/api/api';
import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function EmailVerification() {
    const router = useRouter();
    const searchParams = useSearchParams();

    const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
    const [message, setMessage] = useState('Verifying your email...');

    useEffect(() => {
        const verifyEmail = async () => {
            try {
                const token = searchParams?.get('token');

                if (!token) {
                    setStatus('error');
                    setMessage('Verification token not found.');
                    return;
                }

                const response = await apiRequest({
                    body: { token },
                    method: 'PUT',
                    endpoint: apiConfig.endpoints.auth.emailVerification,
                });

                if (!response.success) {
                    alert(response.message);
                    throw new Error(response.message || 'Failed to verify email');
                }

                setStatus('success');
                setMessage('Email verified successfully! You can now log in.');

            } catch (error) {
                console.error('Verification error:', error);
                setStatus('error');
                setMessage('Error verifying email. The token may have expired or be invalid.');
            }
        };

        verifyEmail();
    }, [searchParams, router]);

    return (
        <div className="min-h-screen flex items-center justify-center p-4">
            <Card className="w-full max-w-md">
                <CardHeader>
                    <CardTitle className="text-center">Email Verification</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    {status === 'loading' && (
                        <div className="flex flex-col items-center space-y-4">
                            <Loader2 className="h-8 w-8 animate-spin" />
                            <p className="text-center">{message}</p>
                        </div>
                    )}

                    {status === 'success' && (
                        <div className="space-y-4">
                            <p className="text-center text-green-600">{message}</p>
                            <Button
                                className="w-full"
                                onClick={() => router.push('/login')}
                            >
                                Go to Login
                            </Button>
                        </div>
                    )}

                    {status === 'error' && (
                        <div className="space-y-4">
                            <p className="text-center text-red-600">{message}</p>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
} 