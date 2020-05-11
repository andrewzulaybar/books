import { Component, AfterViewInit, ElementRef, ViewChild } from '@angular/core';

@Component({
  selector: 'app-navigation-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css'],
})
export class NavigationHeaderComponent implements AfterViewInit {
  @ViewChild('navHeader', { read: ElementRef })
  private navigationHeader: ElementRef;

  public get element(): HTMLElement {
    return this.navigationHeader.nativeElement;
  }

  public ngAfterViewInit() {
    this.addBoxShadow();
    window.addEventListener('scroll', () => this.addBoxShadow());
  }

  private addBoxShadow(): void {
    const opacity: string = this.calculateBoxShadowOpacity().toString();
    this.element.style.setProperty('--box-shadow-opacity', opacity);
  }

  private calculateBoxShadowOpacity(): number {
    const { scrollTop }: { scrollTop: number } = document.documentElement;
    const quarterViewport: number = window.innerHeight / 4;
    return Math.min(scrollTop / quarterViewport, 1);
  }
}
