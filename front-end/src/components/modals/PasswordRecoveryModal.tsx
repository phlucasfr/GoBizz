import { Motion } from "@motionone/solid";
import { validateEmail } from "../../util/Index";
import { createSignal, Show } from "solid-js";
import { sendPasswordRecoveryEmail } from "../../api/api";

interface PasswordRecoveryModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const PasswordRecoveryModal = (props: PasswordRecoveryModalProps) => {
  const [email, setEmail] = createSignal("");
  const [error, setError] = createSignal("");
  const [isSubmitting, setIsSubmitting] = createSignal(false);
  const [successMessage, setSuccessMessage] = createSignal("");

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setSuccessMessage("");

    if (!validateEmail(email())) {
      setError("Por favor, insira um e-mail válido.");
      return;
    }

    setIsSubmitting(true);

    const response = await sendPasswordRecoveryEmail(email());
    setIsSubmitting(false);

    if (response.success) {
      setSuccessMessage(
        "Um e-mail de recuperação foi enviado para o endereço fornecido."
      );
    } else {
      setError(response.message || "Erro ao enviar o e-mail. Tente novamente.");
    }
  };

  return (
    <Show when={props.isOpen}>
      <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <Motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3 }}
          class="w-full max-w-md p-8 bg-white rounded-xl shadow-lg space-y-6 relative"
        >
          <button
            onClick={props.onClose}
            class="absolute top-4 right-4 text-gray-600 hover:text-gray-900"
          >
            ✕
          </button>
          <h1 class="text-2xl font-bold text-center text-indigo-900">
            Recuperação de Senha
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
            <Motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.2 }}
            >
              <label
                for="email"
                class="block text-sm font-medium text-gray-700 mb-1"
              >
                E-mail
              </label>
              <input
                id="email"
                name="email"
                type="email"
                value={email()}
                onInput={(e) => setEmail(e.currentTarget.value)}
                required
                class="block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                placeholder="seu@email.com"
              />
            </Motion.div>

            <Motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.3 }}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <button
                type="submit"
                disabled={isSubmitting()}
                class="w-full bg-gradient-to-r from-indigo-600 to-blue-500 text-white py-3 px-5 rounded-lg shadow-md hover:from-indigo-700 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                {isSubmitting()
                  ? "Enviando..."
                  : "Enviar e-mail de recuperação"}
              </button>
            </Motion.div>
          </form>
        </Motion.div>
      </div>
    </Show>
  );
};

export default PasswordRecoveryModal;
