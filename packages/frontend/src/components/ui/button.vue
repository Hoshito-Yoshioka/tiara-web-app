<script setup lang="ts">
import { Primitive } from 'radix-vue'
import { type ButtonVariants, buttonVariants } from '.'
import { cn } from '@/lib/utils'

interface Props {
  variant?: ButtonVariants['variant']
  size?: ButtonVariants['size']
  asChild?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  asChild: false,
})
</script>

<template>
  <Primitive
    :as-child="asChild"
    :class="cn(buttonVariants({ variant, size }))"
  >
    <slot />
  </Primitive>
</template>
