import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Menu, X } from "lucide-react";

const links = [
  { name: "Vehicle", href: "#" },
  { name: "Control Panel", href: "#" },
  { name: "City Map", href: "#" },
  { name: "Emergency", href: "#" },
  { name: "Traffic Light", href: "#" },
];

export default function Navbar() {
  const [isOpen, setIsOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 50);
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <nav
      className={`fixed top-0 left-0 z-50 w-full transition-all  duration-500 ${
        scrolled ? "bg-black/80 shadow-lg backdrop-blur-md" : "bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900"
      } py-3`}
    >
      <div className="flex justify-between items-center max-w-7xl mx-auto px-6 md:px-12">
        {/* Logo */}
        <motion.div
          className="flex items-center cursor-pointer"
          whileHover={{ scale: 1.1 }}
        >
          <img
            src="/assets/image.png"
            alt="Logo"
            className="shadow-lg border-2 border-white rounded-full w-16 h-16"
          />
          <span className="ml-3 font-mono font-light text-white text-3xl uppercase tracking-wide">
            Code Crafters
          </span>
        </motion.div>

        {/* Desktop Links */}
        <div className="hidden md:flex space-x-10">
          {links.map((link, index) => (
            <motion.a
              key={index}
              href={link.href}
              className="group relative text-gray-300 hover:text-white text-lg transition-all duration-300"
              whileHover={{ scale: 1.1 }}
            >
              {link.name}
              <motion.span
                className="absolute bottom-0 left-0 bg-blue-400 w-full h-0.5 scale-x-0 group-hover:scale-x-100 transition-transform"
              ></motion.span>
            </motion.a>
          ))}
        </div>

        {/* Mobile Menu Button */}
        <button
          className="md:hidden text-white"
          onClick={() => setIsOpen(!isOpen)}
        >
          {isOpen ? <X size={32} /> : <Menu size={32} />}
        </button>
      </div>

      {/* Mobile Menu */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="fixed top-16 left-0 w-full flex flex-col items-center space-y-6 bg-black bg-opacity-90 py-6"
          >
            {links.map((link, index) => (
              <a
                key={index}
                href={link.href}
                className="text-gray-300 hover:text-white text-lg transition duration-300"
                onClick={() => setIsOpen(false)}
              >
                {link.name}
              </a>
            ))}
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  );
}
