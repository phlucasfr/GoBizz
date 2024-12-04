import { Motion } from "@motionone/solid";
import { resetPassword } from "../../api/api";
import { createSignal, Show } from "solid-js";
import { useSearchParams, useNavigate } from "@solidjs/router";

const ResetPasswordPage = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const tokenParam = searchParams.token;
  const token = Array.isArray(tokenParam) ? tokenParam[0] : tokenParam;

  const [error, setError] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [isSubmitting, setIsSubmitting] = createSignal(false);
  const [successMessage, setSuccessMessage] = createSignal("");
  const [confirmPassword, setConfirmPassword] = createSignal("");

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setSuccessMessage("");

    if (!token) {
      setError("Token inválido ou expirado.");
      return;
    }

    if (password() !== confirmPassword()) {
      setError("As senhas não coincidem.");
      return;
    }

    if (password().length < 8) {
      setError("A senha deve ter pelo menos 8 caracteres.");
      return;
    }

    setIsSubmitting(true);

    const response = await resetPassword(token, password());
    setIsSubmitting(false);

    if (response.success) {
      setSuccessMessage(
        "Senha alterada com sucesso! Você já pode fazer login."
      );
    } else {
      setError(
        response.message || "Erro ao redefinir a senha. Tente novamente."
      );
    }
  };

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
          Redefinir Senha
        </h1>

        <Show when={error()}>
          <Motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.3 }}
            class="mb-4 p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg"
          >
            {error()}
          </Motion.div>
        </Show>

        <Show when={successMessage()}>
          <Motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.3 }}
            class="mb-4 p-4 bg-green-100 text-green-800 border-l-4 border-green-500 rounded-lg"
          >
            {successMessage()}
          </Motion.div>
        </Show>

        <form onSubmit={handleSubmit} class="space-y-6">
          <div>
            <label
              for="password"
              class="block text-sm font-medium text-gray-700 mb-1"
            >
              Nova Senha
            </label>
            <input
              id="password"
              name="password"
              type="password"
              value={password()}
              onInput={(e) => setPassword(e.currentTarget.value)}
              required
              class="block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              placeholder="Digite sua nova senha"
            />
          </div>

          <div>
            <label
              for="confirm-password"
              class="block text-sm font-medium text-gray-700 mb-1"
            >
              Confirmar Nova Senha
            </label>
            <input
              id="confirm-password"
              name="confirm-password"
              type="password"
              value={confirmPassword()}
              onInput={(e) => setConfirmPassword(e.currentTarget.value)}
              required
              class="block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              placeholder="Confirme sua nova senha"
            />
          </div>

          <button
            type="submit"
            disabled={isSubmitting()}
            class="w-full bg-gradient-to-r from-indigo-600 to-blue-500 text-white py-3 px-5 rounded-lg shadow-md hover:from-indigo-700 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
          >
            {isSubmitting() ? "Redefinindo..." : "Redefinir Senha"}
          </button>
        </form>

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

export default ResetPasswordPage;
