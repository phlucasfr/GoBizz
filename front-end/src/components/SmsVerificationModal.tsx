import { apiConfig } from "../config/apiConfig";
import { createSignal, Show } from "solid-js";

interface VerificationModalProps {
  onClose: () => void;
  isVisible: boolean;
}

const SmsVerificationModal = (props: VerificationModalProps) => {
  console.log("Renderizando modal", props.isVisible);

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
          body: JSON.stringify({
            id,
            code: code(),
          }),
        }
      );

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || "Falha ao verificar o código.");
      }

      console.log("Código verificado com sucesso!");
      props.onClose();
    } catch (err: any) {
      setError(err.message || "Erro inesperado. Tente novamente.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Show when={props.isVisible}>
      <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
        <div class="w-full max-w-sm p-6 bg-white rounded-lg shadow-lg space-y-6">
          <h2 class="text-2xl font-semibold text-center text-blue-600">
            Verificação de Código
          </h2>
          <p class="text-sm text-center text-gray-500">
            Insira o código de 6 dígitos que enviamos para seu número de
            telefone.
          </p>

          {error() && (
            <div class="p-3 bg-red-100 text-red-800 border-l-4 border-red-500 rounded">
              {error()}
            </div>
          )}

          <form onSubmit={handleSubmit} class="space-y-4">
            <input
              type="text"
              inputMode="numeric"
              maxLength={6}
              placeholder="000000"
              value={code()}
              onInput={handleInput}
              class="w-full text-black text-center text-xl font-medium px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              type="submit"
              disabled={isSubmitting()}
              class={`w-full py-3 px-5 rounded-lg shadow-md text-white ${
                isSubmitting() ? "bg-gray-400" : "bg-blue-600 hover:bg-blue-500"
              } focus:outline-none focus:ring-2 focus:ring-blue-600`}
            >
              {isSubmitting() ? "Verificando..." : "Confirmar"}
            </button>
          </form>

          <button
            onClick={props.onClose}
            class="w-full text-sm text-gray-600 hover:text-gray-800 focus:outline-none"
          >
            Cancelar
          </button>
        </div>
      </div>
    </Show>
  );
};

export default SmsVerificationModal;
