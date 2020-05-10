import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-section-footer',
  templateUrl: './section-footer.component.html',
  styleUrls: ['./section-footer.component.css'],
})
export class SectionFooterComponent implements OnInit {
  @Input() public infoRoute: string;

  constructor() {}

  ngOnInit(): void {}
}
