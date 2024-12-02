import { Motion } from "@motionone/solid";
import { useAuth } from "./context/AuthContext";
import { Home, Info, Mail, Menu, X } from "lucide-solid";
import { createSignal, createMemo, For, Show } from "solid-js";

type SectionKey = "home" | "about" | "contact" | "profile" | "logout";
type Section = {
  id: SectionKey;
  href: string;
  icon: typeof Home;
  title: string;
  action?: () => void;
};

const Header = () => {
  const { isLoggedIn } = useAuth();
  const [isMenuOpen, setIsMenuOpen] = createSignal(false);
  
  const closeMenu = () => setIsMenuOpen(false);

  const sections = createMemo<Section[]>(() => {
    const baseSections: Section[] = [
      {
        id: "about",
        href: "/about",
        icon: Info,
        title: "Sobre Nós",
      },
      {
        id: "contact",
        href: "/contact",
        icon: Mail,
        title: "Contato",
      },
    ];

    return baseSections;
  });

  return (
    <Show when={!isLoggedIn()}>
      <div class="min-h-16">
        <header class="fixed top-0 left-0 right-0 z-50">
          <Motion.div
            initial={{ opacity: 0, y: -50 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            class="bg-gradient-to-r from-blue-600 to-indigo-800 text-white shadow-lg"
          >
            <div class="container mx-auto px-4 sm:px-6 lg:px-8">
              <div class="flex items-center justify-between h-16">
                <div class="flex items-center">
                  <a
                    href="/"
                    class="flex-shrink-0"
                    onClick={() => {
                      setIsMenuOpen(!isMenuOpen());
                    }}
                  >
                    <img class="h-8 w-8" src="/assets/logo.png" alt="Logo" />
                  </a>
                  <div class="hidden md:block">
                    <div class="ml-10 flex items-baseline space-x-4">
                      <For each={sections()}>
                        {(section) => (
                          <Motion.div
                            whileHover={{ scale: 1.05 }}
                            whileTap={{ scale: 0.95 }}
                            class="flex items-center justify-center"
                          >
                            <a
                              href={section.href}
                              onClick={(e) => {
                                if (section.action) {
                                  e.preventDefault();
                                  section.action();
                                }
                              }}
                              class="group relative flex items-center justify-center p-3 rounded-md transition-all duration-300 ease-in-out"
                            >
                              <section.icon
                                class="h-6 w-6 text-gray-300 group-hover:text-white transition-colors duration-300 ease-in-out"
                                aria-hidden="true"
                              />
                              <span class="absolute -bottom-7 text-sm font-medium text-gray-300 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                                {section.title}
                              </span>
                            </a>
                          </Motion.div>
                        )}
                      </For>
                    </div>
                  </div>
                </div>
                <div class="-mr-2 flex md:hidden">
                  <button
                    onClick={() => setIsMenuOpen(!isMenuOpen())}
                    class="bg-indigo-800 inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-indigo-800 focus:ring-white"
                    aria-expanded="false"
                  >
                    <span class="sr-only">Open main menu</span>
                    {isMenuOpen() ? (
                      <X class="block h-6 w-6" aria-hidden="true" />
                    ) : (
                      <Menu class="block h-6 w-6" aria-hidden="true" />
                    )}
                  </button>
                </div>
              </div>
            </div>

            <Motion.div
              class={`${isMenuOpen() ? "block" : "hidden"} md:hidden`}
              initial={{ opacity: 0, y: -20 }}
              animate={{
                opacity: isMenuOpen() ? 1 : 0,
                y: isMenuOpen() ? 0 : -20,
              }}
              transition={{ duration: 0.3 }}
            >
              <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
                <For each={sections()}>
                  {(section) => (
                    <a
                      href={section.href}
                      onClick={(e) => {
                        if (section.action) {
                          e.preventDefault();
                          section.action();
                        }
                        closeMenu();
                      }}
                      class="text-gray-300 hover:bg-indigo-700 hover:text-white px-3 py-2 rounded-md text-base font-medium w-full text-left flex items-center"
                    >
                      <section.icon class="h-5 w-5 mr-2" aria-hidden="true" />
                      {section.title}
                    </a>
                  )}
                </For>
              </div>
            </Motion.div>
          </Motion.div>
        </header>
      </div>
    </Show>
  );
};

export default Header;
