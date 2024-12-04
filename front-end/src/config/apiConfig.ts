export const apiConfig = {
  baseUrl: import.meta.env.VITE_AUTH_SERVICE_API,
  endpoints: {
    company: "/v1/companies",
    sessions: "/v1/sessions",
    recovery: "/v1/companies/recovery",
    companyLogin: "/v1/companies/login",
    resetPassword: "/v1/companies/reset-password",  
    emailVerification: "/v1/companies/email-verification"  
  },
};
