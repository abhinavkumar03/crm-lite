"use client";

import { Toaster } from "sonner";

export default function ToastProvider() {
  return (
    <Toaster
      position="top-right"
      richColors
      expand={false}
      closeButton
      duration={3500}
      toastOptions={{
        style: {
          borderRadius: "18px",
        },
      }}
    />
  );
}