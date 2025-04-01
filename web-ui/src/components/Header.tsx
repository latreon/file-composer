import { FiZap } from 'react-icons/fi'
import Link from 'next/link'

export default function Header() {
    return (
        <header className="py-4 mb-8">
            <div className="flex items-center justify-center">
                <Link href="/" className="flex items-center space-x-2">
                    <div className="bg-primary text-white p-2 rounded-lg">
                        <FiZap className="w-6 h-6" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold">FileCompressor</h1>
                        <p className="text-xs text-gray-500 dark:text-gray-400">Fast & Secure Compression</p>
                    </div>
                </Link>
            </div>
        </header>
    )
} 