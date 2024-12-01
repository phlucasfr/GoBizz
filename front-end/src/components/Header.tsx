import { Motion } from "@motionone/solid";
import { isLoggedIn } from "../App";
import { createSignal, createMemo, For } from "solid-js";

type SectionKey = "home" | "about" | "contact" | "profile" | "logout";
type Section = {
  id: SectionKey;
  href: string;
  icon: string;
  title: string;
};

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = createSignal(false);

  const closeMenu = () => setIsMenuOpen(false);

  const sections = createMemo(() => {
    const baseSections: Section[] = [
      {
        id: "about",
        title: "Sobre Nós",
        icon: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
        href: "/about",
      },
      {
        id: "contact",
        title: "Contato",
        icon: "M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z",
        href: "/contact",
      },
    ];

    if (isLoggedIn()) {
      return [
        ...baseSections,
        {
          id: "profile",
          title: "Perfil",
          icon: "M12 14c-4.418 0-8-3.582-8-8s3.582-8 8-8 8 3.582 8 8-3.582 8-8 8z",
          href: "/profile",
        },
        {
          id: "logout",
          title: "Sair",
          icon: "M16 7v10m2-5H6",
          href: "/logout",
        },
      ];
    }

    return baseSections;
  });

  return (
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
                <a href="/" class="flex-shrink-0">
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
                            class="group relative flex items-center justify-center p-3 rounded-md transition-all duration-300 ease-in-out"
                          >
                            <svg
                              class="h-6 w-6 text-gray-300 group-hover:text-white transition-colors duration-300 ease-in-out"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                              aria-hidden="true"
                            >
                              <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d={section.icon}
                              />
                            </svg>
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
                  <svg
                    class={`${isMenuOpen() ? "hidden" : "block"} h-6 w-6`}
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    aria-hidden="true"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M4 6h16M4 12h16M4 18h16"
                    />
                  </svg>
                  <svg
                    class={`${isMenuOpen() ? "block" : "hidden"} h-6 w-6`}
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    aria-hidden="true"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
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
                    onClick={closeMenu}
                    class="text-gray-300 hover:bg-indigo-700 hover:text-white block px-3 py-2 rounded-md text-base font-medium w-full text-left"
                  >
                    {section.title}
                  </a>
                )}
              </For>
            </div>
          </Motion.div>
        </Motion.div>
      </header>
    </div>
  );
};

export default Header;
