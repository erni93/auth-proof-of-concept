import { Component, OnInit } from '@angular/core';
import { faUserClock } from '@fortawesome/free-solid-svg-icons';
import { faCompass } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit {
  faUserClock = faUserClock;
  faCompass = faCompass;

  constructor() {}

  ngOnInit(): void {}
}
