'use client';

import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Loader2 } from 'lucide-react';
import { apiConfig } from '@/api/config';
import { apiRequest } from '@/api/api';
import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function ResetPasswordContent() {
    const router = useRouter();
    const searchParams = useSearchParams();

    const [token, setToken] = useState('');
    const [error, setError] = useState('');
    const [status, setStatus] = useState<'loading' | 'form' | 'success' | 'error'>('loading');
    const [message, setMessage] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');

    useEffect(() => {
        const tokenFromUrl = searchParams?.get('token');
        if (tokenFromUrl) {
            setToken(tokenFromUrl);
            setStatus('form');
        } else {
            setStatus('error');
            setMessage('Reset token not found.');
        }
    }, [searchParams]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (newPassword !== confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        if (newPassword.length < 8) {
            setError('Password must be at least 8 characters long');
            return;
        }

        setStatus('loading');
        try {
            const response = await apiRequest({
                method: 'PUT',
                endpoint: apiConfig.endpoints.auth.resetPassword,
                body: {
                    token,
                    password: newPassword,
                },
            });

            if (!response.success) {
                throw new Error(response.message || 'Failed to reset password');
            }

            setStatus('success');
            setMessage('Password reset successfully! You can now log in with your new password.');

            setTimeout(() => {
                router.push('/login');
            }, 3000);
        } catch (error) {
            console.error('Reset password error:', error);
            setStatus('error');
            setMessage('Error resetting password. The token may have expired or be invalid.');
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center p-4">
            <Card className="w-full max-w-md">
                <CardHeader>
                    <CardTitle className="text-center">Reset Password</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    {status === 'loading' && (
                        <div className="flex flex-col items-center space-y-4">
                            <Loader2 className="h-8 w-8 animate-spin" />
                            <p className="text-center">{message}</p>
                        </div>
                    )}

                    {status === 'form' && (
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="space-y-2">
                                <Label htmlFor="new-password">New Password</Label>
                                <Input
                                    id="new-password"
                                    type="password"
                                    value={newPassword}
                                    onChange={(e) => setNewPassword(e.target.value)}
                                    required
                                    minLength={8}
                                />
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="confirm-password">Confirm Password</Label>
                                <Input
                                    id="confirm-password"
                                    type="password"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    required
                                    minLength={8}
                                />
                            </div>

                            {error && (
                                <p className="text-sm text-red-500 text-center">
                                    {error}
                                </p>
                            )}

                            <Button type="submit" className="w-full">
                                Reset Password
                            </Button>
                        </form>
                    )}

                    {status === 'success' && (
                        <div className="space-y-4">
                            <p className="text-center text-green-600">{message}</p>
                            <Button className="w-full" onClick={() => router.push('/login')}>
                                Go to Login
                            </Button>
                        </div>
                    )}

                    {status === 'error' && (
                        <div className="space-y-4">
                            <p className="text-center text-red-600">{message}</p>
                            <Button className="w-full" onClick={() => router.push('/login')}>
                                Back to Login
                            </Button>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
