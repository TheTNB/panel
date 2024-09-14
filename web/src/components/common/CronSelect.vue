<script setup lang="ts">
const emit = defineEmits(['update'])

const minutes = [...Array(60).keys()]
const hours = [...Array(24).keys()]
const days = [...Array(31).keys()].map((n) => n + 1)
const months = [...Array(12).keys()].map((n) => n + 1)
const weeks = [...Array(7).keys()].map((n) => n + 1)

const selectedMinuteOption = ref<string>('every')
const selectedMinutes = ref<any>([])
const cycleMinuteStart = ref<number>(0)
const cycleMinuteEnd = ref<number>(0)
const pointMinuteStart = ref<number>(0)
const pointMinuteEnd = ref<number>(0)

const selectedHourOption = ref<string>('every')
const selectedHours = ref<any>([])
const cycleHourStart = ref<number>(0)
const cycleHourEnd = ref<number>(0)
const pointHourStart = ref<number>(0)
const pointHourEnd = ref<number>(0)

const selectedDayOption = ref<string>('every')
const selectedDays = ref<any>([])
const cycleDayStart = ref<number>(1)
const cycleDayEnd = ref<number>(1)
const pointDayStart = ref<number>(1)
const pointDayEnd = ref<number>(1)

const selectedMonthOption = ref<string>('every')
const selectedMonths = ref<any>([])
const cycleMonthStart = ref<number>(1)
const cycleMonthEnd = ref<number>(1)
const pointMonthStart = ref<number>(1)
const pointMonthEnd = ref<number>(1)

const selectedWeekOption = ref<string>('every')
const selectedWeeks = ref<any>([])
const cycleWeekStart = ref<number>(1)
const cycleWeekEnd = ref<number>(1)
const pointWeekStart = ref<number>(1)
const pointWeekEnd = ref<number>(1)

const cronExpression = ref<string>('')

function generateCronExpression() {
  let minuteStr = '*'
  if (selectedMinuteOption.value === 'specific') {
    if (selectedMinutes.value.length == 0) {
      minuteStr = '*'
    } else {
      minuteStr = selectedMinutes.value.join(',')
    }
  } else if (selectedMinuteOption.value === 'cycle') {
    minuteStr = `${cycleMinuteStart.value}-${cycleMinuteEnd.value}`
  } else if (selectedMinuteOption.value === 'point') {
    minuteStr = `${pointMinuteStart.value}/${pointMinuteEnd.value}`
  }

  let hourStr = '*'
  if (selectedHourOption.value === 'specific') {
    if (selectedHours.value.length === 0) {
      hourStr = '*'
    } else {
      hourStr = selectedHours.value.join(',')
    }
  } else if (selectedHourOption.value === 'cycle') {
    hourStr = `${cycleHourStart.value}-${cycleHourEnd.value}`
  } else if (selectedHourOption.value === 'point') {
    hourStr = `${pointHourStart.value}/${pointHourEnd.value}`
  }

  let dayStr = '*'
  if (selectedDayOption.value === 'specific') {
    if (selectedDays.value.length === 0) {
      dayStr = '*'
    } else {
      dayStr = selectedDays.value.join(',')
    }
  } else if (selectedDayOption.value === 'cycle') {
    dayStr = `${cycleDayStart.value}-${cycleDayEnd.value}`
  } else if (selectedDayOption.value === 'point') {
    dayStr = `${pointDayStart.value}/${pointDayEnd.value}`
  }

  let monthStr = '*'
  if (selectedMonthOption.value === 'specific') {
    if (selectedMonths.value.length === 0) {
      monthStr = '*'
    } else {
      monthStr = selectedMonths.value.join(',')
    }
  } else if (selectedMonthOption.value === 'cycle') {
    monthStr = `${cycleMonthStart.value}-${cycleMonthEnd.value}`
  } else if (selectedMonthOption.value === 'point') {
    monthStr = `${pointMonthStart.value}/${pointMonthEnd.value}`
  }

  let weekStr = '*'
  if (selectedWeekOption.value === 'specific') {
    if (selectedWeeks.value.length === 0) {
      weekStr = '*'
    } else {
      weekStr = selectedWeeks.value.join(',')
    }
  } else if (selectedWeekOption.value === 'cycle') {
    weekStr = `${cycleWeekStart.value}-${cycleWeekEnd.value}`
  } else if (selectedWeekOption.value === 'point') {
    weekStr = `${pointWeekStart.value}/${pointWeekEnd.value}`
  }

  cronExpression.value = `${minuteStr} ${hourStr} ${dayStr} ${monthStr} ${weekStr}`
}

watch(
  [
    selectedMinutes,
    selectedHours,
    selectedDays,
    selectedMonths,
    selectedWeeks,
    selectedMinuteOption,
    selectedHourOption,
    selectedDayOption,
    selectedMonthOption,
    selectedWeekOption
  ],
  () => {
    generateCronExpression()
    emit('update', cronExpression.value)
  },
  { deep: true }
)
</script>

<template>
  <n-card>
    <n-tabs type="line" animated>
      <n-tab-pane name="minute" tab="分">
        <n-space vertical>
          <n-radio
            value="every"
            @change="selectedMinuteOption = 'every'"
            :checked="selectedMinuteOption === 'every'"
          >
            每分
          </n-radio>
          <n-space>
            <n-radio
              value="cycle"
              @change="selectedMinuteOption = 'cycle'"
              :checked="selectedMinuteOption === 'cycle'"
            >
              周期
            </n-radio>
            从
            <n-input-number v-model:value="cycleMinuteStart" size="tiny" w-100></n-input-number>
            -
            <n-input-number v-model:value="cycleMinuteEnd" size="tiny" w-100></n-input-number>
            分 (0-59)
          </n-space>
          <n-space>
            <n-radio
              value="point"
              @change="selectedMinuteOption = 'point'"
              :checked="selectedMinuteOption === 'point'"
            >
              按照
            </n-radio>
            从
            <n-input-number v-model:value="pointMinuteStart" size="tiny" w-100></n-input-number>
            分开始，每
            <n-input-number v-model:value="pointMinuteEnd" size="tiny" w-100></n-input-number>
            分执行一次 (0/60)
          </n-space>
          <n-radio
            value="specific"
            @change="selectedMinuteOption = 'specific'"
            :checked="selectedMinuteOption === 'specific'"
          >
            指定
          </n-radio>
          <n-space>
            <n-checkbox-group
              :value="selectedMinutes"
              :disabled="selectedMinuteOption !== 'specific'"
              @update:value="selectedMinutes = $event"
            >
              <n-checkbox v-for="item in minutes" :key="item" :value="item" :label="String(item)" />
            </n-checkbox-group>
          </n-space>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="hour" tab="时">
        <n-space vertical>
          <n-radio
            value="every"
            @change="selectedHourOption = 'every'"
            :checked="selectedHourOption === 'every'"
          >
            每时
          </n-radio>
          <n-space>
            <n-radio
              value="cycle"
              @change="selectedHourOption = 'cycle'"
              :checked="selectedHourOption === 'cycle'"
            >
              周期
            </n-radio>
            从
            <n-input-number v-model:value="cycleHourStart" size="tiny" w-100></n-input-number>
            -
            <n-input-number v-model:value="cycleHourEnd" size="tiny" w-100></n-input-number>
            时 (0-23)
          </n-space>
          <n-space>
            <n-radio
              value="point"
              @change="selectedHourOption = 'point'"
              :checked="selectedHourOption === 'point'"
            >
              按照
            </n-radio>
            从
            <n-input-number v-model:value="pointHourStart" size="tiny" w-100></n-input-number>
            时开始，每
            <n-input-number v-model:value="pointHourEnd" size="tiny" w-100></n-input-number>
            时执行一次 (0/24)
          </n-space>
          <n-radio
            value="specific"
            @change="selectedHourOption = 'specific'"
            :checked="selectedHourOption === 'specific'"
          >
            指定
          </n-radio>
          <n-space>
            <n-checkbox-group
              :value="selectedHours"
              :disabled="selectedHourOption !== 'specific'"
              @update:value="selectedHours = $event"
            >
              <n-checkbox v-for="item in hours" :key="item" :value="item" :label="String(item)" />
            </n-checkbox-group>
          </n-space>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="day" tab="日">
        <n-space vertical>
          <n-radio
            value="every"
            @change="selectedDayOption = 'every'"
            :checked="selectedDayOption === 'every'"
          >
            每日
          </n-radio>
          <n-space>
            <n-radio
              value="cycle"
              @change="selectedDayOption = 'cycle'"
              :checked="selectedDayOption === 'cycle'"
            >
              周期
            </n-radio>
            从
            <n-input-number v-model:value="cycleDayStart" size="tiny" w-100></n-input-number>
            -
            <n-input-number v-model:value="cycleDayEnd" size="tiny" w-100></n-input-number>
            日 (1-31)
          </n-space>
          <n-space>
            <n-radio
              value="point"
              @change="selectedDayOption = 'point'"
              :checked="selectedDayOption === 'point'"
            >
              按照
            </n-radio>
            从
            <n-input-number v-model:value="pointDayStart" size="tiny" w-100></n-input-number>
            日开始，每
            <n-input-number v-model:value="pointDayEnd" size="tiny" w-100></n-input-number>
            日执行一次 (1/31)
          </n-space>
          <n-radio
            value="specific"
            @change="selectedDayOption = 'specific'"
            :checked="selectedDayOption === 'specific'"
          >
            指定
          </n-radio>
          <n-space>
            <n-checkbox-group
              :value="selectedDays"
              :disabled="selectedDayOption !== 'specific'"
              @update:value="selectedDays = $event"
            >
              <n-checkbox v-for="item in days" :key="item" :value="item" :label="String(item)" />
            </n-checkbox-group>
          </n-space>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="month" tab="月">
        <n-space vertical>
          <n-radio
            value="every"
            @change="selectedMonthOption = 'every'"
            :checked="selectedMonthOption === 'every'"
          >
            每月
          </n-radio>
          <n-space>
            <n-radio
              value="cycle"
              @change="selectedMonthOption = 'cycle'"
              :checked="selectedMonthOption === 'cycle'"
            >
              周期
            </n-radio>
            从
            <n-input-number v-model:value="cycleMonthStart" size="tiny" w-100></n-input-number>
            -
            <n-input-number v-model:value="cycleMonthEnd" size="tiny" w-100></n-input-number>
            月 (1-12)
          </n-space>
          <n-space>
            <n-radio
              value="point"
              @change="selectedMonthOption = 'point'"
              :checked="selectedMonthOption === 'point'"
            >
              按照
            </n-radio>
            从
            <n-input-number v-model:value="pointMonthStart" size="tiny" w-100></n-input-number>
            月开始，每
            <n-input-number v-model:value="pointMonthEnd" size="tiny" w-100></n-input-number>
            月执行一次 (1/12)
          </n-space>
          <n-radio
            value="specific"
            @change="selectedMonthOption = 'specific'"
            :checked="selectedMonthOption === 'specific'"
          >
            指定
          </n-radio>
          <n-space>
            <n-checkbox-group
              :value="selectedMonths"
              :disabled="selectedMonthOption !== 'specific'"
              @update:value="selectedMonths = $event"
            >
              <n-checkbox v-for="item in months" :key="item" :value="item" :label="String(item)" />
            </n-checkbox-group>
          </n-space>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="week" tab="周">
        <n-space vertical>
          <n-radio
            value="every"
            @change="selectedWeekOption = 'every'"
            :checked="selectedWeekOption === 'every'"
          >
            每周
          </n-radio>
          <n-space>
            <n-radio
              value="cycle"
              @change="selectedWeekOption = 'cycle'"
              :checked="selectedWeekOption === 'cycle'"
            >
              周期
            </n-radio>
            从
            <n-input-number v-model:value="cycleWeekStart" size="tiny" w-100></n-input-number>
            -
            <n-input-number v-model:value="cycleWeekEnd" size="tiny" w-100></n-input-number>
            周 (1-7)
          </n-space>
          <n-space>
            <n-radio
              value="point"
              @change="selectedWeekOption = 'point'"
              :checked="selectedWeekOption === 'point'"
            >
              按照
            </n-radio>
            第
            <n-input-number v-model:value="pointWeekStart" size="tiny" w-100></n-input-number>
            周的星期
            <n-input-number v-model:value="pointWeekEnd" size="tiny" w-100></n-input-number>
            (1-4 / 1-7)
          </n-space>
          <n-radio
            value="specific"
            @change="selectedWeekOption = 'specific'"
            :checked="selectedWeekOption === 'specific'"
          >
            指定
          </n-radio>
          <n-space>
            <n-checkbox-group
              :value="selectedWeeks"
              :disabled="selectedWeekOption !== 'specific'"
              @update:value="selectedWeeks = $event"
            >
              <n-checkbox v-for="item in weeks" :key="item" :value="item" :label="String(item)" />
            </n-checkbox-group>
          </n-space>
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>
