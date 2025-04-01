export default function Footer() {
    return (
        <footer className="mt-12 py-6 border-t border-gray-200 dark:border-gray-800">
            <div className="text-center text-gray-500 dark:text-gray-400 text-sm">
                <p>
                    &copy; {new Date().getFullYear()} File Compressor - Built with Go and Next.js
                </p>
                <p className="mt-1">
                    <a
                        href="https://github.com/latreon/file-compressor"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-primary hover:underline"
                    >
                        View on GitHub
                    </a>
                </p>
            </div>
        </footer>
    )
} 