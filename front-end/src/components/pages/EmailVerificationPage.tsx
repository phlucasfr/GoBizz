import { Motion } from "@motionone/solid";
import { VerifyCompanyByEmail } from "../../api/api";
import { createSignal, onMount } from "solid-js";
import { useSearchParams, useNavigate } from "@solidjs/router";

const EmailVerificationPage = () => {
  const [searchParams] = useSearchParams();
  const tokenParam = searchParams.token;
  const token = Array.isArray(tokenParam) ? tokenParam[0] : tokenParam;

  const navigate = useNavigate();

  const [error, setError] = createSignal("");
  const [successMessage, setSuccessMessage] = createSignal("");
  const [isVerifying, setIsVerifying] = createSignal(false);

  onMount(async () => {
    if (!token) {
      setError("Token inválido ou ausente.");
      return;
    }

    setIsVerifying(true);

    const response = await VerifyCompanyByEmail(token);

    if (!response.success) {
      setIsVerifying(false);
      return setError(response.message || "Erro ao verificar o e-mail.");
    }

    setSuccessMessage(response.message);
    setIsVerifying(false);
  });

  const handleReturnToMain = () => {
    navigate("/");
  };

  return (
    <div class="flex items-center justify-center h-screen bg-gray-100">
      <Motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        class="w-full max-w-md p-8 bg-white rounded-xl shadow-lg space-y-6"
      >
        <h1 class="text-2xl font-bold text-center text-indigo-900">
          Verificação de E-mail
        </h1>

        {isVerifying() && (
          <p class="text-center text-gray-500">Verificando...</p>
        )}

        {error() && (
          <Motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.3 }}
            class="p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg"
          >
            {error()}
          </Motion.div>
        )}

        {successMessage() && (
          <Motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.3 }}
            class="p-4 bg-green-100 text-green-800 border-l-4 border-green-500 rounded-lg"
          >
            {successMessage()}
          </Motion.div>
        )}

        <button
          onClick={handleReturnToMain}
          class="w-full bg-gray-200 text-gray-800 py-3 px-5 rounded-lg shadow-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-all duration-200"
        >
          Voltar para a Página Principal
        </button>
      </Motion.div>
    </div>
  );
};

export default EmailVerificationPage;
