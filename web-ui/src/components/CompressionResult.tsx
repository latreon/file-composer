import { FiDownload, FiRefreshCw, FiCheckCircle, FiAlertCircle } from 'react-icons/fi'

interface CompressionResultProps {
    result: {
        success: boolean
        message?: string
        downloadLink?: string
        outputSize?: number
        inputSize?: number
    }
    onReset: () => void
}

export default function CompressionResult({ result, onReset }: CompressionResultProps) {
    // Calculate compression percentage if sizes are available
    const compressionPercentage =
        result.inputSize && result.outputSize
            ? 100 - Math.round((result.outputSize / result.inputSize) * 100)
            : null

    // Format file size to human-readable format
    const formatFileSize = (bytes?: number) => {
        if (bytes === undefined) return 'Unknown'

        const units = ['B', 'KB', 'MB', 'GB']
        let size = bytes
        let unitIndex = 0

        while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024
            unitIndex++
        }

        return `${size.toFixed(2)} ${units[unitIndex]}`
    }

    return (
        <div className="card p-6">
            <div className="flex items-center mb-4">
                {result.success ? (
                    <FiCheckCircle className="w-8 h-8 text-green-500 mr-3" />
                ) : (
                    <FiAlertCircle className="w-8 h-8 text-red-500 mr-3" />
                )}
                <h2 className="text-xl font-semibold">
                    {result.success ? 'Compression Successful!' : 'Compression Failed'}
                </h2>
            </div>

            <p className="text-gray-600 dark:text-gray-300 mb-6">
                {result.message || (result.success ? 'Your file has been compressed successfully.' : 'An error occurred during compression.')}
            </p>

            {result.success && compressionPercentage !== null && (
                <div className="bg-gray-100 dark:bg-gray-800 rounded-lg p-4 mb-6">
                    <div className="flex justify-between mb-2">
                        <span className="text-gray-600 dark:text-gray-400">Compression Ratio:</span>
                        <span className="font-semibold text-primary">{compressionPercentage}% smaller</span>
                    </div>

                    <div className="flex justify-between mb-2">
                        <span className="text-gray-600 dark:text-gray-400">Original Size:</span>
                        <span>{formatFileSize(result.inputSize)}</span>
                    </div>

                    <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">Compressed Size:</span>
                        <span>{formatFileSize(result.outputSize)}</span>
                    </div>
                </div>
            )}

            <div className="flex flex-col sm:flex-row gap-3">
                {result.success && result.downloadLink && (
                    <a
                        href={result.downloadLink}
                        className="btn btn-primary flex-1 flex justify-center items-center"
                        download
                    >
                        <FiDownload className="mr-2" />
                        Download Compressed File
                    </a>
                )}

                <button
                    onClick={onReset}
                    className="btn flex-1 border border-gray-300 dark:border-gray-700 flex justify-center items-center"
                >
                    <FiRefreshCw className="mr-2" />
                    Compress Another File
                </button>
            </div>
        </div>
    )
} 