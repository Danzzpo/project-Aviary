<script setup>
import { ref, computed } from 'vue'
import {
  Combobox,
  ComboboxInput,
  ComboboxButton,
  ComboboxOptions,
  ComboboxOption,
  TransitionRoot,
} from '@headlessui/vue'
import { Check, ChevronDown, Search } from 'lucide-vue-next'

const props = defineProps({
  modelValue: [Number, String, Object], // ID yang terpilih
  options: { type: Array, default: () => [] }, // Daftar Burung
  label: String,
  placeholder: String
})

const emit = defineEmits(['update:modelValue'])

const query = ref('')

// Filter Data: Tampilkan 50 teratas yang cocok dengan pencarian agar tidak berat
const filteredOptions = computed(() =>
  query.value === ''
    ? props.options.slice(0, 50) 
    : props.options
        .filter((person) =>
          person.label.toLowerCase().replace(/\s+/g, '').includes(query.value.toLowerCase().replace(/\s+/g, ''))
        )
        .slice(0, 50) 
)

// Helper untuk menampilkan text dari ID yang terpilih
const displayValue = (id) => {
  const selected = props.options.find(o => o.id === id)
  return selected ? selected.label : ''
}
</script>

<template>
  <div class="w-full">
    <label v-if="label" class="block text-sm font-medium text-slate-700 mb-1">{{ label }}</label>
    <Combobox :modelValue="modelValue" @update:modelValue="value => emit('update:modelValue', value)">
      <div class="relative mt-1">
        <div class="relative w-full cursor-default overflow-hidden rounded-lg bg-white text-left border border-slate-200 focus-within:ring-2 focus-within:ring-emerald-500 focus-within:border-emerald-500 sm:text-sm transition shadow-sm">
          <ComboboxInput
            class="w-full border-none py-2.5 pl-10 pr-10 text-sm leading-5 text-slate-900 focus:ring-0 outline-none"
            :displayValue="(id) => displayValue(id)"
            @change="query = $event.target.value"
            :placeholder="placeholder"
            autocomplete="off"
          />
          <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
            <Search :size="16" />
          </div>
          <ComboboxButton class="absolute inset-y-0 right-0 flex items-center pr-2">
            <ChevronDown class="h-5 w-5 text-slate-400" aria-hidden="true" />
          </ComboboxButton>
        </div>
        
        <TransitionRoot
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
          @after-leave="query = ''"
        >
          <ComboboxOptions class="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm z-50 scrollbar-thin scrollbar-thumb-slate-300 scrollbar-track-slate-100">
            <div
              v-if="filteredOptions.length === 0 && query !== ''"
              class="relative cursor-default select-none py-2 px-4 text-slate-700"
            >
              Tidak ditemukan.
            </div>

            <ComboboxOption
              v-for="person in filteredOptions"
              as="template"
              :key="person.id"
              :value="person.id"
              v-slot="{ selected, active }"
            >
              <li
                class="relative cursor-default select-none py-2 pl-10 pr-4 transition"
                :class="{
                  'bg-emerald-50 text-emerald-900': active,
                  'text-slate-900': !active,
                }"
              >
                <span :class="{ 'font-medium': selected, 'font-normal': !selected }" class="block truncate">
                  {{ person.label }}
                  <span class="text-xs text-slate-400 ml-2">({{ person.subLabel }})</span>
                </span>
                <span v-if="selected" class="absolute inset-y-0 left-0 flex items-center pl-3 text-emerald-600">
                  <Check class="h-5 w-5" aria-hidden="true" />
                </span>
              </li>
            </ComboboxOption>
          </ComboboxOptions>
        </TransitionRoot>
      </div>
    </Combobox>
  </div>
</template>