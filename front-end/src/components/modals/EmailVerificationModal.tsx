import { Show } from "solid-js";
import { Motion } from "@motionone/solid";
import { useNavigate } from "@solidjs/router";

interface VerificationModalProps {
  onClose: () => void;
  isVisible: boolean;
}

const EmailVerificationModal = (props: VerificationModalProps) => {
  const navigate = useNavigate();

  const handleReturn = () => {
    props.onClose();
    navigate("/");
  };

  return (
    <Show when={props.isVisible}>
      <Motion.div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Motion.div
          class="bg-white rounded-lg shadow-lg w-full max-w-sm p-6 text-center"
          initial={{ scale: 0.9 }}
          animate={{ scale: 1 }}
          transition={{ duration: 0.3 }}
        >
          <h2 class="text-2xl font-semibold text-indigo-800">
            E-mail Enviado!
          </h2>
          <p class="text-gray-600 mt-4">
            Um e-mail de verificação foi enviado para o endereço cadastrado. Por
            favor, verifique sua caixa de entrada.
          </p>
          <div class="mt-6">
            <button
              onClick={handleReturn}
              class="w-full bg-indigo-600 text-white py-3 rounded-lg font-medium hover:bg-indigo-700 transition"
            >
              Voltar para a Página Principal
            </button>
          </div>
        </Motion.div>
      </Motion.div>
    </Show>
  );
};

export default EmailVerificationModal;
