import { useCallback, useState } from 'react'
import { useDropzone } from 'react-dropzone'
import { FiUpload, FiFile } from 'react-icons/fi'
import { motion } from 'framer-motion'

interface FileUploadProps {
    onFileSelect: (file: File) => void
    isLoading: boolean
    format: string
}

export default function FileUpload({ onFileSelect, isLoading, format }: FileUploadProps) {
    const [fileHover, setFileHover] = useState(false)
    const [selectedFile, setSelectedFile] = useState<File | null>(null)

    const onDrop = useCallback(
        (acceptedFiles: File[]) => {
            const file = acceptedFiles[0]
            if (file) {
                setSelectedFile(file)
                onFileSelect(file)
            }
        },
        [onFileSelect]
    )

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        disabled: isLoading,
        multiple: false,
    })

    return (
        <div
            className={`card border-2 border-dashed transition-all ${isDragActive || fileHover
                ? 'border-primary bg-primary bg-opacity-5'
                : 'border-gray-300 dark:border-gray-700'
                }`}
        >
            <div
                {...getRootProps()}
                className="p-8 text-center cursor-pointer"
                onMouseEnter={() => setFileHover(true)}
                onMouseLeave={() => setFileHover(false)}
            >
                <input {...getInputProps()} />

                {isLoading ? (
                    <div className="py-10 flex flex-col items-center">
                        <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mb-4"></div>
                        <p className="text-gray-600 dark:text-gray-300">Compressing your file...</p>
                    </div>
                ) : (
                    <div className="py-10">
                        <motion.div
                            initial={{ scale: 1 }}
                            animate={{ scale: isDragActive ? 1.1 : 1 }}
                            transition={{ duration: 0.3 }}
                            className="flex justify-center mb-4"
                        >
                            {isDragActive ? (
                                <FiUpload className="w-16 h-16 text-primary" />
                            ) : (
                                <FiFile className="w-16 h-16 text-gray-400 dark:text-gray-500" />
                            )}
                        </motion.div>

                        <div className="space-y-2">
                            <p className="text-lg font-medium">
                                {isDragActive
                                    ? 'Drop your file here'
                                    : selectedFile
                                        ? `${selectedFile.name} selected`
                                        : 'Drag & drop your file here'}
                            </p>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                {selectedFile ? 'Click to select a different file' : 'or click to browse'}
                            </p>
                            <p className="text-xs text-gray-400 dark:text-gray-500 mt-2">
                                Any file type supported • Max size: 100MB
                            </p>
                            <p className="text-xs font-medium text-primary mt-1">
                                Will be compressed as .{format.toLowerCase()}
                            </p>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
} 