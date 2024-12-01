export const apiConfig = {
  baseUrl: import.meta.env.VITE_AUTH_SERVICE_API,
  endpoints: {
    company: "/v1/companies",
    sessions: "/v1/sessions",
    smsVerify: "/v1/companies/sms/verify",
  },
};
