export const apiConfig = {
    baseUrl: process.env.NEXT_PUBLIC_AUTH_SERVICE_API,
    endpoints: {
        auth: {
            logout: "/v1/auth/logout",
            customer: "/v1/auth",
            recovery: "/v1/auth/recovery",
            refreshToken: "/v1/auth/refresh-token",
            customerLogin: "/v1/auth/login",
            resetPassword: "/v1/auth/reset-password",
            validateSession: "/v1/auth/validate-session",
            emailVerification: "/v1/auth/email-verification",
        },

        links: {
            getLink: "/v1/links/:shortUrl",
            createLink: "/v1/links",
            deleteLink: "/v1/links",
            updateLink: "/v1/links/:id",
            getCustomerLinks: "/v1/links/customer/:customerId",
            updateLinkClicks: "/v1/links/:id/clicks",
        },
    },
} as const; 