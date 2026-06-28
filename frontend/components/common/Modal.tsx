"use client";

interface ModalProps {
    open: boolean;
    title: string;
    onClose: () => void;
    children: React.ReactNode;
}

export default function Modal({
    open,
    title,
    onClose,
    children,
}: ModalProps) {

    if (!open) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">

            <div className="w-full max-w-lg rounded-lg bg-white shadow-lg">

                <div className="flex items-center justify-between border-b p-4">

                    <h2 className="text-lg font-semibold">
                        {title}
                    </h2>

                    <button onClick={onClose}>
                        ✕
                    </button>

                </div>

                <div className="p-6">

                    {children}

                </div>

            </div>

        </div>
    );
}