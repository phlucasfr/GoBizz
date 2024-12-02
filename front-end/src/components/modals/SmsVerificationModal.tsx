import { Motion } from "@motionone/solid";
import { apiConfig } from "../../config/apiConfig";
import { useNavigate } from "@solidjs/router";
import { validateSession } from "../../api/api";
import { createSignal, Show } from "solid-js";

interface VerificationModalProps {
  onClose: () => void;
  isVisible: boolean;
}

const SmsVerificationModal = (props: VerificationModalProps) => {
  const navigate = useNavigate();
  const [code, setCode] = createSignal("");
  const [error, setError] = createSignal("");
  const [isSubmitting, setIsSubmitting] = createSignal(false);

  const handleInput = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const value = target.value.replace(/\D/g, "");
    setCode(value.slice(0, 6));
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      if (code().length !== 6) {
        throw new Error("Por favor, insira o código completo de 6 dígitos.");
      }

      const id = localStorage.getItem("id");
      if (!id) {
        throw new Error("ID não encontrado. Tente novamente.");
      }

      const response = await fetch(
        `${apiConfig.baseUrl}${apiConfig.endpoints.smsVerify}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ id, code: code() }),
          credentials: "include",
        }
      );
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Falha ao verificar o código.");
      }

      const sessionData = await validateSession();
      if (sessionData instanceof Error) {
        throw new Error(sessionData.message || "Falha na verificação da sessão.");
      }

      if (sessionData.isValid) {
        navigate(`/home/${id}`);
      } else {
        navigate("/login");
      }

      props.onClose();
    } catch (err: any) {
      setError(err.message || "Erro inesperado. Tente novamente.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setCode("");
    setError("");
    setIsSubmitting(false);
    props.onClose();
  };

  return (
    <Show when={props.isVisible}>
      <Motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.3 }}
        class="fixed inset-0 z-50 flex items-center justify-center bg-white"
      >
        <div class="w-full max-w-md p-6 space-y-8">
          <Motion.div
            initial={{ scale: 0.9, y: 20 }}
            animate={{ scale: 1, y: 0 }}
            transition={{ duration: 0.3 }}
          >
            <h2 class="text-3xl font-bold text-center text-indigo-900">Verificação de Código</h2>
            <p class="mt-4 text-lg text-center text-gray-600">
              Insira o código de 6 dígitos que enviamos para seu número de telefone.
            </p>
          </Motion.div>

          <Show when={error()}>
            <Motion.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3 }}
              class="p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg"
              role="alert"
            >
              {error()}
            </Motion.div>
          </Show>

          <form onSubmit={handleSubmit} class="space-y-6">
            <Motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: 0.1 }}
            >
              <label for="verification-code" class="sr-only">Código de verificação</label>
              <input
                id="verification-code"
                type="text"
                inputMode="numeric"
                maxLength={6}
                placeholder="000000"
                value={code()}
                onInput={handleInput}
                class="w-full text-center text-3xl font-medium px-5 py-4 border-2 border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                aria-label="Insira o código de 6 dígitos"
              />
            </Motion.div>
            <Motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: 0.2 }}
            >
              <button
                type="submit"
                disabled={isSubmitting()}
                class="w-full bg-gradient-to-r from-indigo-600 to-blue-500 text-white py-4 px-6 rounded-lg text-lg font-semibold shadow-md hover:from-indigo-700 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                {isSubmitting() ? "Verificando..." : "Confirmar"}
              </button>
            </Motion.div>
          </form>

          <Motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3, delay: 0.3 }}
          >
            <button
              onClick={handleClose}
              class="w-full text-lg text-indigo-600 hover:text-indigo-500 focus:outline-none transition-colors duration-200"
            >
              Cancelar
            </button>
          </Motion.div>
        </div>
      </Motion.div>
    </Show>
  );
};

export default SmsVerificationModal;

