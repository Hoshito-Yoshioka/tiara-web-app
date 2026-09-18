import { describe, it, expect } from 'vitest'
import { HOME_TITLE, BRAND_TITLE, PAGE_META } from '@/lib/seo'

/** コードポイント基準の文字数（サロゲートペアを 1 文字と数える） */
function countCodePoints(str: string): number {
  return [...str].length
}

describe('SEO コピーの機械検証', () => {
  describe('HOME_TITLE', () => {
    it('必須キーワード（函館・ニュークラブ・キャバクラ・ティアラ・TIARA・【公式】）をすべて含む', () => {
      const requiredKeywords = [
        '函館',
        'ニュークラブ',
        'キャバクラ',
        'ティアラ',
        'TIARA',
        '【公式】',
      ]
      for (const keyword of requiredKeywords) {
        expect(HOME_TITLE, `HOME_TITLE に「${keyword}」が含まれること`).toContain(keyword)
      }
    })
  })

  describe('BRAND_TITLE', () => {
    it('「函館」「ニュークラブ」を含む', () => {
      expect(BRAND_TITLE).toContain('函館')
      expect(BRAND_TITLE).toContain('ニュークラブ')
    })
  })

  describe('PAGE_META の description', () => {
    it('全ページの description が 90〜125 文字（コードポイント基準）である', () => {
      for (const [page, meta] of Object.entries(PAGE_META)) {
        const length = countCodePoints(meta.description)
        expect(
          length,
          `${page} の description は 90〜125 文字であること（実際: ${length} 文字）`
        ).toBeGreaterThanOrEqual(90)
        expect(
          length,
          `${page} の description は 90〜125 文字であること（実際: ${length} 文字）`
        ).toBeLessThanOrEqual(125)
      }
    })

    it('全ページの description が「函館」を含む', () => {
      for (const [page, meta] of Object.entries(PAGE_META)) {
        expect(meta.description, `${page} の description に「函館」が含まれること`).toContain(
          '函館'
        )
      }
    })
  })

  describe('PAGE_META の下層 title', () => {
    it('下層ページの title 固有部が互いに重複しない', () => {
      const titles = Object.values(PAGE_META)
        .map((meta) => meta.title)
        .filter((title): title is string => title !== undefined)
      expect(titles.length, '下層ページの title 固有部が定義されていること').toBeGreaterThan(0)
      expect(new Set(titles).size, '下層ページの title 固有部が一意であること').toBe(titles.length)
    })
  })
})
