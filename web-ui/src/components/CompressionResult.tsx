import { FiDownload, FiRefreshCw } from 'react-icons/fi'

interface CompressionResultProps {
    result: {
        success: boolean
        message?: string
        downloadLink?: string
        outputSize?: number
        inputSize?: number
    }
    format?: string
    onReset: () => void
}

export default function CompressionResult({ result, format, onReset }: CompressionResultProps) {
    const compressionRatio = result.inputSize && result.outputSize
        ? ((1 - result.outputSize / result.inputSize) * 100).toFixed(1)
        : null

    const formatBytes = (bytes?: number) => {
        if (!bytes) return '0 B'
        const sizes = ['B', 'KB', 'MB', 'GB']
        const i = Math.floor(Math.log(bytes) / Math.log(1024))
        return `${(bytes / Math.pow(1024, i)).toFixed(2)} ${sizes[i]}`
    }

    return (
        <div className="card p-6">
            {result.success ? (
                <>
                    <div className="text-center mb-6">
                        <h2 className="text-2xl font-bold text-green-600 dark:text-green-400 mb-2">
                            {format === 'pdf' ? 'PDF Optimized Successfully!' : 'Compression Complete!'}
                        </h2>
                        {compressionRatio && (
                            <p className="text-gray-600 dark:text-gray-300">
                                {format === 'pdf' ? 'Reduced' : 'Compressed'} by {compressionRatio}%
                            </p>
                        )}
                    </div>

                    <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 mb-6">
                        <div className="grid grid-cols-2 gap-4">
                            <div className="text-center">
                                <p className="text-sm text-gray-500 dark:text-gray-400">Original Size</p>
                                <p className="text-lg font-semibold">{formatBytes(result.inputSize)}</p>
                            </div>
                            <div className="text-center">
                                <p className="text-sm text-gray-500 dark:text-gray-400">
                                    {format === 'pdf' ? 'Optimized Size' : 'Compressed Size'}
                                </p>
                                <p className="text-lg font-semibold">{formatBytes(result.outputSize)}</p>
                            </div>
                        </div>
                    </div>

                    <div className="flex flex-col sm:flex-row gap-4">
                        <a
                            href={result.downloadLink}
                            className="flex-1 btn btn-primary flex items-center justify-center gap-2"
                            download
                        >
                            <FiDownload />
                            <span>Download {format === 'pdf' ? 'Optimized PDF' : 'Compressed File'}</span>
                        </a>
                        <button
                            onClick={onReset}
                            className="flex-1 btn btn-secondary flex items-center justify-center gap-2"
                        >
                            <FiRefreshCw />
                            <span>Compress Another File</span>
                        </button>
                    </div>
                </>
            ) : (
                <div className="text-center">
                    <h2 className="text-xl font-bold text-red-600 dark:text-red-400 mb-4">
                        {format === 'pdf' ? 'PDF Optimization Failed' : 'Compression Failed'}
                    </h2>
                    <p className="text-gray-600 dark:text-gray-300 mb-6">{result.message}</p>
                    <button
                        onClick={onReset}
                        className="btn btn-primary flex items-center justify-center gap-2 mx-auto"
                    >
                        <FiRefreshCw />
                        <span>Try Again</span>
                    </button>
                </div>
            )}
        </div>
    )
} 