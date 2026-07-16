"use client";

import Image from "next/image";
import Link from "next/link";
import { Facebook, Instagram, Twitter } from "@mynaui/icons-react";

const links = [
  { label: "Home", href: "#home" },
  { label: "About", href: "#about" },
  { label: "Menu", href: "#menu" },
  { label: "Gallery", href: "#gallery" },
  { label: "Contact", href: "#contact" },
];

export default function Footer() {
  return (
    <footer className="pt-24 lg:pt-36">
      <div className="mx-auto max-w-7xl px-6 pb-8 lg:px-10">
        <div className="mt-32 grid gap-14 text-center md:grid-cols-3 md:text-left">
          {/* Quick Links */}

          <div>
            <h3 className="font-display text-3xl text-primary">Quick Links</h3>

            <ul className="mt-6 space-y-3">
              {links.map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-text-secondary transition hover:text-primary"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Brand */}

          <div className="flex flex-col items-center">
            <Image
              src="/assets/gericht.png"
              alt="Gericht"
              width={180}
              height={60}
            />

            <p className="mt-6 max-w-sm text-center leading-8 text-text-secondary">
              Great food is more than a meal—it is an experience shared,
              remembered, and celebrated with the people who matter most.
            </p>

            <Image
              src="/assets/spoon.png"
              alt=""
              width={48}
              height={48}
              className="mt-6"
            />

            <div className="mt-8 flex gap-5">
              <Link
                href="#"
                className="rounded-full text-primary border border-border p-3 transition hover:border-primary hover:text-primary"
              >
                <Facebook size={18} />
              </Link>

              <Link
                href="#"
                className="rounded-full text-primary border border-border p-3 transition hover:border-primary hover:text-primary"
              >
                <Instagram size={18} />
              </Link>

              <Link
                href="#"
                className="rounded-full text-primary border border-border p-3 transition hover:border-primary hover:text-primary"
              >
                <Twitter size={18} />
              </Link>
            </div>
          </div>

          {/* Hours */}

          <div className="md:text-right">
            <h3 className="font-display text-3xl text-primary">
              Opening Hours
            </h3>

            <div className="mt-6 space-y-5 text-text-secondary">
              <div>
                <p className="font-medium text-white">Monday – Friday</p>

                <p>11:00 AM – 11:00 PM</p>
              </div>

              <div>
                <p className="font-medium text-white">Saturday – Sunday</p>

                <p>10:00 AM – Midnight</p>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-16 border-t border-border pt-6 text-center text-sm text-text-muted">
          © {new Date().getFullYear()} Gericht Restaurant. Crafted with passion.
        </div>
      </div>
    </footer>
  );
}
