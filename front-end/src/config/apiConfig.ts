export const apiConfig = {
  baseUrl: import.meta.env.AUTH_SERVICE_API
  ,
  endpoints: {
    company: "/v1/companies",
    smsVerify: "/v1/companies/sms/verify",
  },
};
