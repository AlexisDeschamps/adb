<template>
  <adb-page :title="'Voter Info Form - Phone Banking'">
    <form action id="canvassingForm" autocomplete="off">
      <fieldset :disabled="loading">
        <div class="form-group">
          <label for="notes"> Notes </label>
          <textarea
            id="notes"
            type="text"
            class="form-control"
            v-model="supporter.notes"
            :style="resizeNotesStyle"
            @focus="resizeNotes"
            ref="notes"
          />
        </div>

        <br />
        <p>Either email or phone number is required.</p>

        <div class="row">
          <div class="form-group col-sm-6">
            <label for="email">
              Email
            </label>
            <input id="email" class="form-control" v-model="supporter.email" placeholder="Email" />
          </div>
          <div class="form-group col-sm-6">
            <label for="phone">
              Phone
            </label>
            <input id="phone" class="form-control" v-model="supporter.phone" placeholder="Phone" />
          </div>
        </div>

        <div class="row">
          <div class="form-group col-sm-6">
            <label for="firstName">
              First Name
            </label>
            <input
              id="firstName"
              class="form-control"
              v-model="supporter.firstName"
              placeholder="First Name"
            />
          </div>
          <div class="form-group col-sm-6">
            <label for="lastName">
              Last Name
            </label>
            <input
              id="lastName"
              class="form-control"
              v-model="supporter.lastName"
              placeholder="Last Name"
            />
          </div>
        </div>

        <div class="form-group">
          <input type="checkbox" id="isSupporter" v-model="supporter.isSupporter" />
          <label for="isSupporter">
            Is Supporter (e.g. they picked up the phone and want to be included in future
            communications)
          </label>
        </div>
        <div class="form-group">
          <input type="checkbox" id="voter" v-model="supporter.voter" />
          <label for="voter"> Eligible Berkeley Voter </label>
        </div>

        <div class="form-group">
          <label>Location</label>
          <canvassing-address
            v-on:change="onAddressChange"
            :initialAddress1="supporter.locationAddress1"
            :initialAddress2="supporter.locationAddress2"
            :initialCity="supporter.locationCity"
            :initialState="supporter.locationState"
            :initialZip="supporter.locationZip"
          />
        </div>

        <div class="form-group">
          <label for="source">Source</label>
          <select id="source" class="form-control" v-model="supporter.source">
            <option value="phone">Phone Banking</option>
            <option value="canvass">In Person Canvassing</option>
            <option value="canvass">House Party</option>
          </select>
        </div>
        <div class="form-group">
          <label for="dateSourced">Date Sourced</label>
          <input
            id="dateSourced"
            class="form-control"
            type="date"
            v-model="supporter.dateSourced"
          />
        </div>
        <br />

        <div>
          <h3>Requests</h3>
          <div class="form-group">
            <input type="checkbox" id="requestedLawnSign" v-model="supporter.requestedLawnSign" />
            <label for="requestedLawnSign"> Requested Lawn Sign </label>
          </div>
          <div class="form-group">
            <input type="checkbox" id="requestedPoster" v-model="supporter.requestedPoster" />
            <label for="requestedPoster"> Requested Poster </label>
          </div>
        </div>

        <div>
          <h3>Issues</h3>
          <div class="form-group">
            <input type="checkbox" id="issueHousing" v-model="supporter.issueHousing" />
            <label for="issueHousing"> Housing </label>
          </div>
          <div class="form-group">
            <input type="checkbox" id="issueHomelessness" v-model="supporter.issueHomelessness" />
            <label for="issueHomelessness"> Homelessness </label>
          </div>
          <div class="form-group">
            <input type="checkbox" id="issueClimate" v-model="supporter.issueClimate" />
            <label for="issueClimate"> Climate </label>
          </div>
          <div class="form-group">
            <input type="checkbox" id="issuePublicSafety" v-model="supporter.issuePublicSafety" />
            <label for="issuePublicSafety"> Public Safety </label>
          </div>
          <div class="form-group">
            <input
              type="checkbox"
              id="issuePoliceAccountability"
              v-model="supporter.issuePoliceAccountability"
            />
            <label for="issuePoliceAccountability"> Police Accountability </label>
          </div>
          <div class="form-group">
            <input type="checkbox" id="issueTransit" v-model="supporter.issueTransit" />
            <label for="issueTransit"> Transit </label>
          </div>
          <div class="form-group">
            <input
              type="checkbox"
              id="issueEconomicEquality"
              v-model="supporter.issueEconomicEquality"
            />
            <label for="issueEconomicEquality">Economic Equality</label>
          </div>
          <div class="form-group">
            <input type="checkbox" id="issuePublicHealth" v-model="supporter.issuePublicHealth" />
            <label for="issuePublicHealth"> Public Health </label>
          </div>
          <div class="form-group">
            <input type="checkbox" id="issueAnimalRights" v-model="supporter.issueAnimalRights" />
            <label for="issueAnimalRights"> Animal Rights </label>
          </div>
        </div>

        <div>
          <h3>Support</h3>

          <div class="form-group">
            <input type="checkbox" id="interestDonate" v-model="supporter.interestDonate" />
            <label for="interestDonate"> Donating </label>
          </div>

          <div class="form-group">
            <input
              type="checkbox"
              id="interestAttendEvent"
              v-model="supporter.interestAttendEvent"
            />
            <label for="interestAttendEvent"> Attend Event </label>
          </div>

          <div class="form-group">
            <input type="checkbox" id="interestVolunteer" v-model="supporter.interestVolunteer" />
            <label for="interestVolunteer"> Volunteer </label>
          </div>

          <div class="form-group">
            <input type="checkbox" id="interestHostEvent" v-model="supporter.interestHostEvent" />
            <label for="interestHostEvent"> Host Event (e.g. house party) </label>
          </div>
        </div>

        <div class="form-group">
          <input type="checkbox" id="requiresFollowup" v-model="supporter.requiresFollowup" />
          <label for="requiresFollowup">
            Requires Followup (for something not covered by another option)
          </label>
        </div>

        <!--
        <label for="canvasser"> Canvasser </label>
        <input id="canvasser" type="text" class="form-control" v-model="supporter.canvasser" />

        <label for="canvassLeader"> Canvass Leader </label>
        <input id="canvassLeader" type="text" class="form-control" v-model="supporter.canvassLeader" />
        -->
      </fieldset>
    </form>
    <center>
      <button
        class="btn btn-success btn-lg"
        id="submit-button"
        v-on:click="save"
        :disabled="saving"
      >
        <span>Save</span>
      </button>
    </center>
  </adb-page>
</template>

<script lang="ts">
import Vue from 'vue';
import AdbPage from './AdbPage.vue';
import CanvassingAddress from './components/CanvassingAddress.vue';
import { setFlashMessageSuccessCookie, flashMessage } from './flash_message';

function getDateTodayStr() {
  var d = new Date();
  var year = '' + d.getFullYear();
  var rawMonth = d.getMonth() + 1;
  var month = rawMonth > 9 ? '' + rawMonth : '0' + rawMonth;
  var rawDate = d.getDate();
  var date = rawDate > 9 ? '' + rawDate : '0' + rawDate;
  var validDateString = year + '-' + month + '-' + date;
  return validDateString;
}

function emptySupporter() {
  return {
    id: 0,
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
    locationAddress1: '',
    locationAddress2: '',
    locationCity: '',
    locationState: '',
    locationZip: '',
    source: 'phone',
    dateSourced: getDateTodayStr(),
    requestedLawnSign: false,
    requestedPoster: false,
    voter: true,

    issueHousing: false,
    issueHomelessness: false,
    issueClimate: false,
    issuePublicSafety: false,
    issuePoliceAccountability: false,
    issueTransit: false,
    issueEconomicEquality: false,
    issuePublicHealth: false,
    issueAnimalRights: false,

    interestDonate: false,
    interestAttendEvent: false,
    interestVolunteer: false,
    interestHostEvent: false,

    requiresFollowup: false,

    notes: '',
  };
}

function snakeToCamelCase(str: string): string {
  let newStr = '';
  let capitalizeNext = false;
  for (let c of str) {
    if (c === '_') {
      capitalizeNext = true;
      continue;
    }
    if (capitalizeNext) {
      capitalizeNext = false;
      newStr += c.toUpperCase();
      continue;
    }
    newStr += c;
  }
  return newStr;
}

function camelCaseToSnake(str: string): string {
  let newStr = '';
  for (let c of str) {
    if (/[A-Z]/.test(c)) {
      newStr += '_' + c.toLowerCase();
    } else {
      newStr += c;
    }
  }
  return newStr;
}

function assignSupporterJSONToSupporter(supporter: any, supporterJSON: any) {
  for (let jsonField in supporterJSON) {
    let val = supporterJSON[jsonField];
    if (jsonField == 'date_sourced') {
      // date from the server looks like 2020-03-15T13:00:00Z, cut off
      // everything before "T" so it fits our dates.
      val = val.split('T', 1)[0];
    }
    supporter[snakeToCamelCase(jsonField)] = val;
  }
}

function supporterToJSON(supporter: any): string {
  let ret: any = {};
  for (let field in supporter) {
    let val = supporter[field];
    if (typeof val === 'string') {
      val = val.trim();
    }
    ret[camelCaseToSnake(field)] = val;
  }
  return JSON.stringify(ret);
}

export default Vue.extend({
  components: {
    AdbPage,
    CanvassingAddress,
  },
  props: {
    id: String,
    newOnSave: Boolean,
  },
  data() {
    return {
      loading: false,
      saving: false,

      notesHeight: 'auto',

      supporter: emptySupporter(),
      oldSupporter: emptySupporter(),
    };
  },
  methods: {
    // TODO: use correct type for 'e'
    onAddressChange(e: any) {
      this.supporter.locationAddress1 = e.address1;
      this.supporter.locationAddress2 = e.address2;
      this.supporter.locationCity = e.city;
      this.supporter.locationState = e.state;
      this.supporter.locationZip = e.zip;
    },

    updateSupporter() {
      if (Number(this.supporter.id) == 0) {
        return;
      }

      this.loading = true;
      $.ajax({
        url: '/canvass/supporter/get/' + this.supporter.id,
        method: 'GET',
        dataType: 'json',
        success: (data) => {
          this.loading = false;
          if (data.status === 'error') {
            flashMessage('Error : ' + data.message, true);
            return;
          }

          // data.status === "success"
          this.supporter = emptySupporter();
          this.oldSupporter = emptySupporter();
          assignSupporterJSONToSupporter(this.supporter, data.supporter);
          assignSupporterJSONToSupporter(this.oldSupporter, data.supporter);
        },
        error: () => {
          this.loading = false;
          flashMessage('Server error, could not get data.', true);
        },
      });
    },

    save() {
      if (this.supporter.phone.trim() === '' && this.supporter.email.trim() === '') {
        flashMessage('Error: At least one of phone or email must be set.', true);
        return;
      }
      this.saving = true;
      $.ajax({
        url: '/canvass/supporter/save',
        method: 'POST',
        contentType: 'application/json',
        data: supporterToJSON(this.supporter),
        success: (data) => {
          this.saving = false;
          let parsed = JSON.parse(data);
          if (parsed.status === 'error') {
            flashMessage('Error: ' + parsed.message, true);
            return;
          }

          if (this.newOnSave) {
            flashMessage('Saved!', false);
            this.supporter = emptySupporter();
            this.oldSupporter = emptySupporter();
            window.scrollTo(0, 0);
            return;
          }

          if (parsed.redirect) {
            setFlashMessageSuccessCookie('Saved!');
            window.location.href = parsed.redirect;
            return; // unreachable
          }

          flashMessage('Saved!', false);
          // Re-fetch supporter from database just in case it's changed.
          this.updateSupporter();
        },
        error: () => {
          this.saving = false;
          flashMessage('Server error, did not save data.', true);
        },
      });
    },
    resizeNotes() {
      this.notesHeight = 'auto';
      this.$nextTick(() => {
        let notesEl = this.$refs.notes as HTMLElement;
        console.log('scrollHeight: ', notesEl.scrollHeight);
        if (notesEl.scrollHeight >= 54) {
          this.notesHeight = notesEl.scrollHeight + 'px';
        }
      });
    },
    dirty() {
      return JSON.stringify(this.supporter) != JSON.stringify(this.oldSupporter);
    },
  },
  created() {
    if (Number(this.id) != 0) {
      this.supporter.id = Number(this.id);
      this.updateSupporter();
    }

    window.addEventListener('beforeunload', (e) => {
      if (!this.dirty()) {
        return;
      }

      console.log('User data looks dirty');
      console.log('oldSupporter: ', JSON.stringify(this.oldSupporter));
      console.log('supporter: ', JSON.stringify(this.supporter));

      e.preventDefault();
      // MDN says returnValue is unused, but still required by Chrome.
      // https://developer.mozilla.org/en-US/docs/Web/Events/beforeunload
      e.returnValue = '';
    });

    $(window).bind('resize', () => {
      this.resizeNotes();
    });
  },
  computed: {
    resizeNotesStyle: function(): Object {
      return {
        resize: 'none',
        height: this.notesHeight,
        overflow: 'hidden',
      };
    },
  },
  watch: {
    'supporter.notes': function() {
      this.resizeNotes();
    },
  },
});
</script>
