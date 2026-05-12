import { describe, it, expect } from 'vitest'
import { cn } from '@/lib/utils'

describe('cn ユーティリティ', () => {
  it('複数のクラスをマージする', () => {
    const result = cn('text-white', 'bg-black')
    expect(result).toContain('text-white')
    expect(result).toContain('bg-black')
  })

  it('Tailwind の衝突するクラスを後勝ちでマージする', () => {
    const result = cn('text-red-500', 'text-blue-500')
    expect(result).toBe('text-blue-500')
  })

  it('falsy な値を無視する', () => {
    const result = cn('text-white', undefined, null, false, 'bg-black')
    expect(result).toContain('text-white')
    expect(result).toContain('bg-black')
  })
})
