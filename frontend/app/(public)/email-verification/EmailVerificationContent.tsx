'use client';

import { Button } from '@/components/ui/button';
import { Loader2 } from 'lucide-react';
import { apiConfig } from '@/api/config';
import { apiRequest } from '@/api/api';
import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

function EmailVerificationContent() {
    const router = useRouter();
    const searchParams = useSearchParams();

    const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
    const [message, setMessage] = useState('Verificando seu e-mail...');

    useEffect(() => {
        const verifyEmail = async () => {
            try {
                const token = searchParams?.get('token');

                if (!token) {
                    setStatus('error');
                    setMessage('Token de verificação não encontrado.');
                    return;
                }

                const response = await apiRequest({
                    body: { token },
                    method: 'PUT',
                    endpoint: apiConfig.endpoints.auth.emailVerification,
                });

                if (!response.success) {
                    alert(response.message);
                    throw new Error(response.message || 'Falha na verificação do e-mail');
                }

                setStatus('success');
                setMessage('E-mail verificado com sucesso! Agora você pode fazer login.');
            } catch (error) {
                console.error('Erro na verificação:', error);
                setStatus('error');
                setMessage('Erro na verificação do e-mail. O token pode ter expirado ou ser inválido.');
            }
        };

        verifyEmail();
    }, [searchParams]);

    return (
        <div className="min-h-screen flex items-center justify-center p-4">
            <Card className="w-full max-w-md">
                <CardHeader>
                    <CardTitle className="text-center">Verificação de E-mail</CardTitle>
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
                            <Button className="w-full" onClick={() => router.push('/login')}>
                                Ir para Login
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

export default EmailVerificationContent;
